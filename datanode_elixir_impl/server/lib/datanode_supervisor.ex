defmodule DataNode.Supervisor do
  use Supervisor

  def start_link(init_arg) do
    Supervisor.start_link(__MODULE__, init_arg, name: __MODULE__)
  end

  @impl true
  def init(_args) do
    host = System.get_env("NAME_NODE_HOST")
    port = System.get_env("NAME_NODE_PORT")

    IO.puts("HOST: #{host} and PORT: #{port}")

    children = [
      {DataNodeServer, {host, port}}
    ]

    Supervisor.init(children, strategy: :one_for_one)
  end
end
