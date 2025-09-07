defmodule DataNode.Server do
  alias ElixirLS.LanguageServer.Providers.Completion.Reducers.Struct
  use GenServer

  def start_link(data) do
    GenServer.start(__MODULE__, data, name: __MODULE__)
  end

  def register_request() do
    GenServer.call(__MODULE__, :register)
  end

  def heartbeat() do
    GenServer.cast(__MODULE__, :heartbeat)
  end

  def request() do
    GenServer.cast(__MODULE__, :request)
  end

  def stop() do
    GenServer.stop(__MODULE__, :normal)
  end

  # Server callbacks.

  @impl true
  def init({host, port}) do
    opts = [:binary, {:packet, 0}, {:active, false}]

    case :gen_tcp.connect(to_charlist(host), String.to_integer(port), opts) do
      {:ok, socket} ->
        IO.puts("Connected to the server successfully")
        # Now we need to send the registration request so that our node would be registered.
        GenServer.call(__MODULE__, :register)
        {:ok, socket}

      {:error, reason} ->
        IO.puts("Failed to connect to the server, #{inspect(reason)}")
        {:stop, reason}
    end
  end

  @impl true
  def handle_call(:register, _from, socket) do
    {:ok, {ip, port}} = :inet.sockname(socket)

    node_info = %DataNode.Struct.Node{
      req_type: "REGISTER",
      name: System.get_env("NAME"),
      port: port
    }

    json_node_info = JSON.encode!(node_info)

    case :gen_tcp.send(socket, json_node_info) do
      :ok ->
        IO.puts("Node registered successfully")
        # Now We need to send the heartbeat and accepts the further requests.
        send(self(), :heartbeat)
        send(self(), :req)
        {:reply, :ok, socket}

      {:error, reason} ->
        IO.puts("Failed to register the node, #{inspect(reason)}")
        # Crashing the process.
        raise "Registration Failed: #{inspect(reason)}"
    end
  end

  @impl true
  def handle_cast(:heartbeat, socket) do
    heartbeat = %DataNode.Struct.HeartBeat{
      req_type: "HEARTBEAT",
      node_name: "datanode-1",
      # Timestamp
      timestamp: DateTime.utc_now() |> DateTime.to_unix()
    }

    json_encoded_data = JSON.encode!(heartbeat)

    case :gen_tcp.send(socket, json_encoded_data) do
      :ok ->
        Process.send_after(self(), :heartbeat, 3000)
        {:noreply, socket}

      {:error, reason} ->
        IO.inspect("Failed to send the heartbeat #{inspect(reason)}")
    end
  end

  @impl true
  def handle_cast(:request, socket) do
    case :gen_tcp.recv(socket, 0) do
      {:ok, "GET\n"} ->
        IO.puts("GET request recieved")

        case :gen_tcp.recv(socket, 0) do
          {:ok, data} ->
            case JSON.decode(data) do
              {:ok, decoded_data} ->
                get_request_data = struct(DataNode.Struct.GetRequest, decoded_data)

                case handle_get_request(get_request_data) do
                  {:ok, data} ->
                    response = %DataNode.Struct.Response{
                      type: "SUCCESS",
                      message: binary_data
                    }

                    case :gen_tcp.send(socket, JSON.encode!(response)) do
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

                    case :gen_tcp.send(socket, JSON.encode!(response)) do
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

                :gen_tcp.send(socket, JSON.encode!(response))
            end

          {:error, reason} ->
            IO.puts("Failed to recieve the get request data, #{inspect(reason)}")
        end

      {:ok, "POST\n"} ->
        IO.puts("POST request recieved")

        case :gen_tcp.recv(socket, 0) do
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

                    case :gen_tcp.send(socket, JSON.encode!(response)) do
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

                    case :gen_tcp.send(socket, JSON.encode!(response)) do
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

    send(self(), :req)
    {:noreply, socket}
  end

  defp handle_get_request(data) do
    filename = data.file_name
    chunkid = data.chunk_id

    file_path = Path.join("storage", filename, chunkid)

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
          end

        {:error, reason} ->
          IO.puts("Failed to create the directory, #{inspect(reason)}")
          {:error, reason}
      end
    end
  end

  @impl true
  def handle_info(msg, socket) do
    case msg do
      :req ->
        handle_cast(msg, socket)

      :heartbeat ->
        handle_cast(msg, socket)
    end
  end

  @impl true
  def terminate(reason, socket) do
    IO.puts("Terminating: #{inspect(reason)}")
    :gen_tcp.close(socket)
    :ok
  end
end
