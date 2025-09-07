defmodule DataNode.Struct.Node do
  @derive JSON.Encoder
  defstruct [:name, :port] # Right now only taking the node name and the port at which it is running.
end
