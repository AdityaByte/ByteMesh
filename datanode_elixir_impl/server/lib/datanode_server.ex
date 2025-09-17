defmodule DataNode.Server do
  use GenServer

  def start_link(port \\ 0) do
    GenServer.start(__MODULE__, port, name: __MODULE__)
  end

  def get_server_addr() do
    GenServer.call(__MODULE__, :get_addr)
  end

  def stop() do
    GenServer.stop(__MODULE__, :normal)
  end

  # Server Callbacks
  @impl true
  def init(port) do
    opts = [:binary, {:packet, :line}, {:reuseaddr, true}]
    IO.puts("Starting server...")

    case :gen_tcp.listen(port, opts) do
      {:ok, listen_socket} ->
        {:ok, {host, port}} = :inet.sockname(listen_socket)
        host = :inet.ntoa(host) |> to_string()
        IO.puts("#{System.get_env("NAME")} is listening to #{host}:#{port}")

        # Starting the accept connection process seperately.
        Task.start(fn -> accept_conn(listen_socket) end)

        {:ok, %{host: host, port: port, listen_socket: listen_socket}}

      {:error, reason} ->
        IO.puts("Failed to listen at the desired port: #{inspect(reason)}")
        {:stop, reason}
    end
  end

  defp accept_conn(listen_socket) do
    case :gen_tcp.accept(listen_socket) do
      {:ok, socket} ->
        {:ok, {host, port}} = :inet.sockname(socket)

        IO.puts(
          "Client has been connected successfully, client info: Host: #{host} and port: #{port}"
        )

        # Now we have to handle each and every connection in a seperate process.
        Task.start(fn -> handle_conn(socket) end)

      {:error, reason} ->
        IO.puts("Failed to connect to the client, #{inspect(reason)}")
    end

    # Keeps on listening.
    accept_conn(listen_socket)
  end

  defp handle_conn(socket) do
    case :gen_tcp.recv(socket, 0) do
      {:ok, "GET\n"} ->
        IO.puts("GET request recieved")

        case :gen_tcp.recv(socket, :line) do
          {:ok, data} ->
            case JSON.decode(data) do
              {:ok, decoded_data} ->
                get_request_data = struct(DataNode.Struct.GetRequest, decoded_data)

                case handle_get_request(get_request_data) do
                  {:ok, data} ->
                    response = %DataNode.Struct.Response{
                      type: "SUCCESS",
                      message: data
                    }

                    case :gen_tcp.send(socket, JSON.encode!(response) <> "\n") do
                      :ok ->
                        IO.puts("Get request fulfilled successfully")

                      {:error, reason} ->
                        IO.puts("Failed to send the get success response, #{inspect(reason)}")
                    end

                  {:error, reason} ->
                    IO.puts(inspect(reason))

                    response = %DataNode.Struct.Response{
                      type: "FAILED",
                      message: inspect(reason)
                    }

                    case :gen_tcp.send(socket, JSON.encode!(response) <> "\n") do
                      :ok ->
                        IO.puts("Failed response sent successfully of get request")

                      {:error, reason} ->
                        IO.puts(
                          "Failed to send the failed response of get request, #{inspect(reason)}"
                        )
                    end
                end

              {:error, reason} ->
                IO.puts("Failed to decode the get request data, #{inspect(reason)}")

                response = %DataNode.Struct.Response{
                  type: "FAILED",
                  message: "Invalid JSON"
                }

                :gen_tcp.send(socket, JSON.encode!(response) <> "\n")
            end

          {:error, reason} ->
            IO.puts("Failed to recieve the get request data, #{inspect(reason)}")
        end

      {:ok, "POST\n"} ->
        IO.puts("POST request recieved")

        case :gen_tcp.recv(socket, :line) do
          {:ok, data} ->
            case JSON.decode(data) do
              {:ok, decoded_data} ->
                chunk_data = struct(DataNode.Struct.Chunk, decoded_data)

                case handle_post_request(chunk_data) do
                  :ok ->
                    response = %DataNode.Struct.Response{
                      type: "SUCCESS",
                      message: "Chunk Saved successfully to the node #{System.get_env("NAME")}"
                    }

                    case :gen_tcp.send(socket, JSON.encode!(response) <> "\n") do
                      :ok ->
                        IO.puts("POST request response sent successfully")

                      {:error, reason} ->
                        IO.puts("Failed to send the POST request response, #{inspect(reason)}")
                    end

                  {:error, reason} ->
                    response = %DataNode.Struct.Response{
                      type: "FAILED",
                      message: inspect(reason)
                    }

                    case :gen_tcp.send(socket, JSON.encode!(response) <> "\n") do
                      :ok ->
                        IO.puts("ERROR POST request response sent successfully")

                      {:error, reason} ->
                        IO.puts(
                          "Failed to send the errorfull post request response, #{inspect(reason)}"
                        )
                    end
                end

              {:error, reason} ->
                IO.puts("Failed to decode the JSON, #{inspect(reason)}")

                response = %DataNode.Struct.Response{
                  type: "FAILED",
                  message: "Invalid JSON"
                }

                :gen_tcp.send(socket, JSON.encode!(response))
            end

          {:error, reason} ->
            IO.puts("Failed to recieve the data, #{inspect(reason)}")
        end

      {:error, reason} ->
        IO.puts("Recv error: #{inspect(reason)}")
    end

    handle_conn(socket)
  end

  defp handle_get_request(data) do
    filename = data.file_name
    chunkid = data.chunk_id

    file_path = Path.join(["storage", filename, chunkid])

    case File.read(file_path) do
      {:ok, data} ->
        {:ok, data}

      {:error, reason} ->
        {:error, "ERROR: Failed to read the file: #{inspect(reason)}"}
    end
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

      case File.mkdir_p(new_dir_path) do
        :ok ->
          IO.puts("Directory created successfully.")

          case File.write(temp_path, data, [:binary]) do
            :ok ->
              IO.puts("Data written successfully to the temporary file.")
              # When the data has been successfully written to the temporary file
              # we need to rename the file.
              case File.rename(temp_path, final_path) do
                :ok ->
                  IO.puts("File renamed successfully")
                  :ok

                {:error, reason} ->
                  IO.puts("Failed to rename the file #{inspect(reason)}")
                  {:error, reason}
              end

            {:error, reason} ->
              IO.puts("Failed to write the temporary file")
              {:error, reason}
          end

        {:error, reason} ->
          IO.puts("Failed to create the directory, #{inspect(reason)}")
          {:error, reason}
      end
    end
  end

  @impl true
  def handle_call(:get_addr, _from, state) do
    {:reply, {state.host, state.port}, state}
  end

  @impl true
  def terminate(reason, data) do
    %{listen_socket: listen_socket} = data
    :gen_tcp.close(listen_socket)
  end
end
