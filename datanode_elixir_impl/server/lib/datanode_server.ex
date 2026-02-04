defmodule DataNode.Server do
  use GenServer

  def start_link(port \\ 0) do
    GenServer.start(__MODULE__, port, name: __MODULE__)
  end

  def get_server_addr() do
    IO.puts("get server addr is called..")
    GenServer.call(__MODULE__, :get_addr)
  end

  def stop() do
    GenServer.stop(__MODULE__, :normal)
  end

  # Server Callbacks
  @impl true
  def init(port) do
    IO.puts("Starting datanode server...")
    opts = [:binary, {:packet, :line}, {:reuseaddr, true}]

    {:ok, listen_socket} = :gen_tcp.listen(port, opts)
    {:ok, {host, port}} = :inet.sockname(listen_socket)
    host = :inet.ntoa(host) |> to_string()

    IO.puts("#{System.get_env("NAME")} is listening to #{host}:#{port}")
    spawn_link(fn -> accept_loop(listen_socket) end)

    {:ok, %{host: host, port: port, listen_socket: listen_socket}}
  end

  defp accept_loop(listen_socket) do
    {:ok, socket} = :gen_tcp.accept(listen_socket)
    :ok = :inet.setopts(socket, [{:active, false}, {:packet, :line}])
    IO.puts("Accepted connection successfully..")
    {:ok, {host, port}} = :inet.peername(socket)
    host = :inet.ntoa(host) |> to_string()
    IO.puts(
      "Client has been connected successfully, client info: Host: #{host} and port: #{port}"
    )
    DataNode.Connection.start_link(socket)
    accept_loop(listen_socket)
  end

  @impl true
  def handle_call(:get_addr, _from, state) do
    IO.puts(
      "fetching the server host and port in the handle_call function of datanodeserver, #{state.host} and #{state.port}"
    )

    {:reply, {state.host, state.port}, state}
  end

  @impl true
  def terminate(_reason, data) do
    %{listen_socket: listen_socket} = data
    :gen_tcp.close(listen_socket)
  end
end
