# List of examples
file(GLOB_RECURSE SRCS *.c)

# List of multithreaded examples 
SET (SRCSTHREAD run_encap_decap.c run_bls.c)

# Add the binary tree directory to the search path for linking and include files
link_directories (${PROJECT_BINARY_DIR}/src
                  /usr/local/lib)

include_directories (${PROJECT_SOURCE_DIR}/include
                     /usr/local/include)

# Examples that are not multi-threaded
foreach(example ${SRCS})
  # Extract the filename without an extension
  get_filename_component(target ${example} NAME_WE)

  add_executable(${target} ${example})

  target_link_libraries(${target} pqnist)
endforeach(example)

# Examples that are multi-threaded
if(${CMAKE_SYSTEM_NAME} MATCHES "Linux")
  foreach(example ${SRCSTHREAD})
    # Extract the filename without an extension
    get_filename_component(base ${example} NAME_WE)
    set(target ${base}_threads)

    add_executable(${target} ${example})

    target_compile_options(${target} PUBLIC -pthread -fopenmp)

    target_link_libraries(${target} pqnist -fopenmp)
  endforeach(example)
endif(${CMAKE_SYSTEM_NAME} MATCHES "Linux")
