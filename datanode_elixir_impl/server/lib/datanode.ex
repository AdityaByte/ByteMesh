defmodule DataNode.Server do
  alias ElixirLS.LanguageServer.Providers.Completion.Reducers.Struct
  use GenServer

  def start_link(data) do
    GenServer.start(__MODULE__, data, name: __MODULE__)
  end

  def register_request() do
    GenServer.call(__MODULE__, :register)
  end

  def heartbeat() do
    GenServer.cast(__MODULE__, :heartbeat)
  end

  def request() do
    GenServer.cast(__MODULE__, :request)
  end

  def stop() do
    GenServer.stop(__MODULE__, :normal)
  end

  # Server callbacks.

  @impl true
  def init({host, port}) do
    opts = [:binary, {:packet, 0}, {:active, false}]
    :gen_tcp.connect(to_charlist(host), String.to_integer(port), opts)
  end

  @impl true
  def handle_call(:register, _from, socket) do
    {:ok, {ip, port}} = :inet.sockname(socket)
    node_info = %DataNode.Struct.Node{
      req_type: "REGISTER",
      name: System.get_env("NAME"),
      port: port
    }
    json_node_info = JSON.encode!(node_info)
    case :gen_tcp.send(socket, json_node_info) do
      :ok ->
        IO.puts("Node registered successfully")
        {:reply, :ok, socket}
      {:error, reason} ->
        IO.puts("Failed to register the node, #{inspect(reason)}")
        # Crashing the process.
        raise "Registration Failed: #{inspect(reason)}"
    end
  end

  @impl true
  def handle_cast(:heartbeat, socket) do
    heartbeat = %DataNode.Struct.HeartBeat{
      req_type: "HEARTBEAT",
      node_name: "datanode-1",
      timestamp: DateTime.utc_now() |> DateTime.to_unix() # Timestamp
    }

    json_encoded_data = JSON.encode!(heartbeat)
    case :gen_tcp.send(socket, json_encoded_data) do
      :ok ->
        Process.send_after(self(), :heartbeat, 3000)
        {:noreply, socket}
      {:error, reason} ->
        IO.inspect("Failed to send the heartbeat #{inspect(reason)}")
    end

  end

  @impl true
  def handle_cast(:request, socket) do

    case :gen_tcp.recv(socket, 0) do
      {:ok, "GET\n"} ->
        IO.puts("GET request recieved")
        {:ok, data} = :gen_tcp.recv(socket, 0)
        # Now I need to decode the json data to my struct.
        decoded_data = JSON.decode(data)
        chunk_data = struct(DataNode.Struct.Chunk, decoded_data)
        handle_get_request(chunk_data)
        {:reply, "GET", socket}
      {:ok, "POST\n"} ->
        IO.puts("POST request recieved")
        {:reply, "POST", socket}
      {:error, reason} ->
        IO.puts("Recv error: #{inspect(reason)}")
        {:reply, {:error, reason}, socket}
    end

    send(self(), :req)
    {:noreply, socket}
  end

  defp handle_get_request(chunk_data) do
    # get 
  end

  @impl true
  def handle_info(msg, socket) do
    case msg do
      :req ->
        handle_cast(msg, socket)
      :heartbeat ->
        handle_cast(msg, socket)
    end
  end

  @impl true
  def terminate(reason, socket) do
    IO.puts("Terminating: #{inspect(reason)}")
    :gen_tcp.close(socket)
    :ok
  end

end
