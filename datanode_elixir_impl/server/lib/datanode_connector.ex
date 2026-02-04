defmodule DataNode.Connector do
  @moduledoc """
  DataNode.Connector is responsible for connecting to the name node and registering the datanode with the name node.
  It also sends the heartbeat to the name node every 30 seconds.
  """
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
    {:ok, nil, {:continue, {:connect, host, port}}}
  end

  def handle_continue({:connect, host, port}, _state) do
    opts = [:binary, {:packet, 0}, {:active, false}]
    {:ok, socket} = :gen_tcp.connect(to_charlist(host), String.to_integer(port), opts)
    IO.puts("Connected to the server successfully..")
    # Now we need to fetch the server addr immediately.
    # Do the registration after some time maybe after 5 seconds.
    send(self(), :register)
    {:noreply, socket}
  end

  @impl true
  def handle_info(:register, socket) do
    # Firstly i need to send the HEALTH Verb that the request is of the health request.
    IO.puts("Sending the register request.")
    :ok = :gen_tcp.send(socket, "REGISTER" <> "\n")

    {host, port} = DataNode.Server.get_server_addr()
    IO.puts("Fetched host and port from datanode server: #{host}:#{port}")

    node_info = %DataNode.Struct.Node{
      name: System.get_env("NAME"),
      host: host,
      port: port,
      time_stamp: DateTime.utc_now() |> DateTime.to_unix()
    }

    IO.puts("Encoding the register node data.")
    json_node_info = JSON.encode!(node_info) <> "\n"

    :ok = :gen_tcp.send(socket, json_node_info)

    IO.puts("Node registered successfully")
    # Now We need to send the heartbeat and accepts the further requests.
    Process.send_after(self(), :heartbeat, 30_000)
    {:noreply, socket}
  end

  @impl true
  def handle_info(:heartbeat, socket) do
    IO.puts("Sending heartbeat")

    :ok = :gen_tcp.send(socket, "HEARTBEAT" <> "\n")
    heartbeat = %DataNode.Struct.HeartBeat{
      node_name: System.get_env("NAME"),
      # Timestamp
      timestamp: DateTime.utc_now() |> DateTime.to_unix()
    }
    json_encoded_data = JSON.encode!(heartbeat) <> "\n"
    :ok = :gen_tcp.send(socket, json_encoded_data)
    Process.send_after(self(), :heartbeat, 30_000)
    {:noreply, socket}
  end

  @impl true
  def terminate(reason, socket) do
    IO.puts("Terminating: #{inspect(reason)}")
    :gen_tcp.close(socket)
  end
end
