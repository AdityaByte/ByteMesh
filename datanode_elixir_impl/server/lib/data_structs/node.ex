defmodule DataNode.struct.Node do
  @derive JSON.Encoder
  defstruct [:req_type, :name, :port] # Right now only taking the node name and the port at which it is running.
end
