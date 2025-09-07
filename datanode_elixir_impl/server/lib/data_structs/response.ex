defmodule DataNode.Struct.Response do
  @derive JSON.Encoder
  defstruct [
    :type,
    :message
  ]
end
