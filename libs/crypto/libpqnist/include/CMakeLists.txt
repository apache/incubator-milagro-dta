# List of headers
file(GLOB headers "amcl/*.h")

install(FILES ${headers}
        DESTINATION ${INSTALL_INCLUDESUBDIR})

# Copy header to golang binary directory
file(COPY ${headers} DESTINATION ${PROJECT_BINARY_DIR}/go/amcl)
      
