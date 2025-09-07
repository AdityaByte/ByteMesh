defmodule DataNode.Struct.HeartBeat do
  @derive JSON.Encoder
  defstruct [:node_name, :timestamp]
end
