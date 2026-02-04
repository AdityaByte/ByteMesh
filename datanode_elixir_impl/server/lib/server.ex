defmodule Server do
  use Application

  @impl true
  def start(_type, _args) do
    Envy.auto_load()

    host = System.get_env("NAME_NODE_HOST")
    port = System.get_env("NAME_NODE_PORT")

    IO.puts("Connecting to the namenode server running on #{host}:#{port}")

    children = [
      {DataNode.Server, 0},
      {DataNode.Connector, {host, port}}
    ]

    Supervisor.start_link(children, strategy: :one_for_one, name: __MODULE__)
  end

  def stop(_state) do
    Supervisor.stop(__MODULE__, :normal)
  end
end
