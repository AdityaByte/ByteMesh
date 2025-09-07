defmodule DataNode.Util do
  def is_all_field_present?(struct) do
    struct
      |> Map.from_struct()
      |> Map.values()
      |> Enum.all?(fn val -> val not in [nil, ""] end)
  end
end
