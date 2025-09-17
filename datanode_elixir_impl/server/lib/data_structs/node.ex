defmodule DataNode.Struct.Node do
  @derive JSON.Encoder
  # Right now only taking the node name and the port at which it is running.
  defstruct [:name, :host, :port, :time_stamp]
end
