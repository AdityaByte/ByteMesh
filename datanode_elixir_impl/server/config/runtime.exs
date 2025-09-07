import Config

config :server, :name_node_url,
  System.get_env("NAME_NODE_HOST")
  System.get_env("NAME_NODE_PORT")
  System.get_env("NAME")
