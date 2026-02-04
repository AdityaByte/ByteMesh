defmodule DataNode.Connection do
  @moduledoc """
  DataNode.Connection handles the connection request, GET and POST
  """
  use GenServer

  def start_link(socket) do
    IO.puts(inspect(socket))
    GenServer.start(__MODULE__, socket, name: __MODULE__)
  end

  # Connection callbacks.
  @impl true
  def init(socket) do
    send(self(), :handle_connection)
    {:ok, socket}
  end

  @impl true
  def handle_info(:handle_connection, socket) do
    # Now we need to accept the data from the socket.
    IO.puts("The socket value in connection is, #{inspect(socket)}")
    data = :gen_tcp.recv(socket, 0)
    IO.inspect(data)

    case data do
      {:ok, req_verb} ->
        case String.trim(req_verb) do
          "GET" ->
            IO.puts("GET request received.")
            {:ok, data} = :gen_tcp.recv(socket, 0)
            data = String.trim(data)
            {:ok, decoded_data} = JSON.decode(data)
            get_request_data = struct(DataNode.Struct.GetRequest, decoded_data)

            {:ok, data} = handle_get_request(get_request_data)

            response = %DataNode.Struct.Response{
              type: "SUCCESS",
              message: data
            }

            :ok = :gen_tcp.send(socket, JSON.encode!(response) <> "\n")

            send(self(), :handle_connection)
            {:noreply, socket}

          "POST" ->
            IO.puts("Post request received.")
            {:ok, data} = :gen_tcp.recv(socket, 0)
            data = String.trim(data)
            {:ok, decoded_data} = JSON.decode(data)

            chunk_data = struct(DataNode.Struct.Chunk, decoded_data)
            :ok = handle_post_request(chunk_data)

            response = %DataNode.Struct.Response{
              type: "SUCCESS",
              message: "Chunk Saved successfully to the node #{System.get_env("NAME")}"
            }

            :ok = :gen_tcp.send(socket, JSON.encode!(response) <> "\n")

            send(self(), :handle_connection)
            {:noreply, socket}

          other ->
            IO.puts("Request verb recieved-> #{other}")
            {:stop, :unknown_verb, socket}

        end
      {:error, :closed} ->
        IO.puts("Client closed connection")
        {:stop, :normal, socket}
      {:error, reason} ->
        IO.puts("recv error: #{inspect(reason)}")
        {:stop, reason, socket}
    end
  end

  defp handle_get_request(data) do
    filename = data.file_name
    chunkid = data.chunk_id

    file_path = Path.join(["storage", filename, chunkid])

    File.read(file_path)
  end

  defp handle_post_request(chunk) do
    if !DataNode.Util.is_all_field_present?(chunk) do
      {:error, "Invalid Chunk Data, Some fields are not present."}
    else
      # Else we need to get the filename and fileid and the chunkdata.
      filename = chunk.file_name
      # usually as chunk1, chunkn
      fileid = chunk.file_id
      data = chunk.data

      temp_path = Path.join("storage/#{filename}", "#{fileid}.tmp")
      final_path = Path.join("storage/#{filename}", fileid)

      new_dir_path = Path.join("storage", filename)

      :ok = File.mkdir_p(new_dir_path)
      IO.puts("Directory created successfully...")
      :ok = File.write(temp_path, data, [:binary])
      :ok = File.rename(temp_path, final_path)
    end
  end
end
