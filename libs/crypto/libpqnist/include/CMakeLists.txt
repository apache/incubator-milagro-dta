# List of headers
file(GLOB headers "pqnist/*.h")

install(FILES ${headers}
        DESTINATION ${INSTALL_INCLUDESUBDIR})
