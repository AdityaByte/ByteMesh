defmodule DataNode.Struct.HeartBeat do
  @derive JSON.Encoder
  defstruct [:req_type, :node_name, :timestamp]
end
