# Post request chunk data schema.

defmodule DataNode.Struct.Chunk do
  @derive JSON.Encoder
  defstruct [
    :file_name,
    :file_id,
    :data
  ]
end
