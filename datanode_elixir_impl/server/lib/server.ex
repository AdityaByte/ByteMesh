defmodule Server.Application do
  use Application

  @impl true
  def start(_type, _args) do
    # Loading the environment variable.
    Envy.auto_load()

    children = [
      DataNode.Supervisor
    ]

    opts = [strategy: :one_for_one, name: DataNode.Supervisor]
    Supervisor.init(children, opts)
  end
end
