defmodule DataNode.Connector do
  use GenServer

  def start_link(data) do
    GenServer.start(__MODULE__, data, name: __MODULE__)
  end

  def register_request() do
    GenServer.cast(__MODULE__, :register)
  end

  def heartbeat() do
    GenServer.cast(__MODULE__, :heartbeat)
  end

  def stop() do
    GenServer.stop(__MODULE__, :normal)
  end

  # Connector callbacks.

  @impl true
  def init({host, port}) do
    opts = [:binary, {:packet, 0}, {:active, false}]

    case :gen_tcp.connect(to_charlist(host), String.to_integer(port), opts) do
      {:ok, socket} ->
        IO.puts("Connected to the server successfully")
        # Now we need to send the registration request so that our node would be registered.
        send(self(), :register)
        {:ok, socket}

      {:error, reason} ->
        IO.puts("Failed to connect to the server, #{inspect(reason)}")
        {:stop, reason}
    end
  end

  @impl true
  def handle_info(:register, socket) do
    # Firstly i need to send the HEALTH Verb that the request is of the health request.
    case :gen_tcp.send(socket, "REGISTER" <> "\n") do
      :ok ->
        {host, port} = DataNode.Server.get_server_addr()
        IO.puts("Fetched host and port from datanode server: #{host}:#{port}")

        node_info = %DataNode.Struct.Node{
          name: System.get_env("NAME"),
          host: host,
          port: port,
          time_stamp: DateTime.utc_now() |> DateTime.to_unix()
        }

        json_node_info = JSON.encode!(node_info) <> "\n"

        case :gen_tcp.send(socket, json_node_info) do
          :ok ->
            IO.puts("Node registered successfully")
            # Now We need to send the heartbeat and accepts the further requests.
            Process.send_after(self(), :heartbeat, 30_000)
            {:noreply, socket}

          {:error, reason} ->
            IO.puts("Failed to register the node, #{inspect(reason)}")
            {:noreply, socket}
        end

      {:error, reason} ->
        IO.puts("Failed to send the Register request, #{inspect(reason)}")
        {:noreply, socket}
    end
  end

  @impl true
  def handle_info(:heartbeat, socket) do
    IO.puts("Sending heartbeat")

    case :gen_tcp.send(socket, "HEARTBEAT" <> "\n") do
      :ok ->
        heartbeat = %DataNode.Struct.HeartBeat{
          node_name: System.get_env("NAME"),
          # Timestamp
          timestamp: DateTime.utc_now() |> DateTime.to_unix()
        }

        json_encoded_data = JSON.encode!(heartbeat) <> "\n"

        case :gen_tcp.send(socket, json_encoded_data) do
          :ok ->
            Process.send_after(self(), :heartbeat, 30_000)
            {:noreply, socket}

          {:error, reason} ->
            IO.inspect("Failed to send the heartbeat #{inspect(reason)}")
        end

      {:error, reason} ->
        IO.puts("Failed to send the Health request, #{inspect(reason)}")
    end

    {:noreply, socket}
  end

  @impl true
  def terminate(reason, socket) do
    IO.puts("Terminating: #{inspect(reason)}")
    :gen_tcp.close(socket)
  end
end
