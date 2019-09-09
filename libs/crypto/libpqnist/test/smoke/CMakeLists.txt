# List of tests
file(GLOB_RECURSE SRCS *.c)

# List of multithreaded examples 
SET (SRCSTHREAD test_encap_decap.c test_bls.c)

# Add the binary tree directory to the search path for linking and include files
link_directories (${PROJECT_BINARY_DIR}/src
                  /usr/local/lib)

include_directories (${PROJECT_SOURCE_DIR}/include
                     /usr/local/include)

# define macro to simplify adding tests
macro(do_test arg result)
  add_test(${arg} ${TARGET_SYSTEM_EMULATOR} ${arg}${CMAKE_EXECUTABLE_SUFFIX})
  set_tests_properties(${arg} PROPERTIES PASS_REGULAR_EXPRESSION ${result})
endmacro()

foreach(test ${SRCS})
  # Extract the filename without an extension
  get_filename_component(target ${test} NAME_WE)

  add_executable(${target} ${test})

  target_link_libraries(${target} pqnist)

  do_test(${target} "SUCCESS")
endforeach(test)

# Multi threading tests
if(${CMAKE_SYSTEM_NAME} MATCHES "Linux")
  foreach(test ${SRCSTHREAD})
    # Extract the filename without an extension
    get_filename_component(base ${test} NAME_WE)
    set(target ${base}_threads)

    add_executable(${target} ${test})

    target_compile_options(${target} PUBLIC -fopenmp)

    target_link_libraries(${target} pqnist -fopenmp)

    do_test(${target} "SUCCESS")
  endforeach(test)
endif(${CMAKE_SYSTEM_NAME} MATCHES "Linux")
