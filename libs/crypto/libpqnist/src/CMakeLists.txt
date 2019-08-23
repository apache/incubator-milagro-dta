file(GLOB_RECURSE SOURCES *.c)

set(target "pqnist")

link_directories(${CMAKE_CURRENT_BINARY_DIR}
                 /usr/local/lib)

include_directories (${PROJECT_SOURCE_DIR}/include
                     /usr/local/include)

add_library(${target} ${LIB_TYPE} ${SOURCES})

target_link_libraries (${target} amcl_bls_BLS381 amcl_pairing_BLS381 amcl_curve_BLS381 amcl_core oqs)

set_target_properties(${target}
                      PROPERTIES VERSION
                                 ${BUILD_VERSION}
                                 SOVERSION
                                 ${VERSION_MAJOR})

install(TARGETS ${target} DESTINATION lib PERMISSIONS
        OWNER_WRITE OWNER_READ OWNER_EXECUTE
        GROUP_READ GROUP_EXECUTE
        WORLD_READ WORLD_EXECUTE)

