defmodule DataNode.Struct.GetRequest do
  @derive JSON.Encoder
  defstruct [
    :file_name,
    :chunk_id
  ]
end
