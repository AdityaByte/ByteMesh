import Config

config :data_node_1, :name_node_url,
  System.get_env("NAME_NODE_HOST") || "localhost"
  System.get_env("NAME_NODE_PORT") || 8080
  System.get_env("NAME")
