import sys

"""
Generates a docker-compose file with n clients
Usage: python3 generate_dc.py n
Where n is the number of clients to generate
"""


def main(n):
    file = open("docker-compose-dev.yaml", "w")
    file.write('version: "3.9"\n')
    file.write("name: tp0\n")
    file.write("services:\n")
    server(file)
    file.write("\n")
    for i in range(int(n)):
        client(file, i + 1)
        file.write("\n")
    networks(file)


def server(file):
    """
    Writes the server section of the docker-compose file
    """
    file.write("  server:\n")
    file.write("    container_name: server\n")
    file.write("    image: server:latest\n")
    file.write("    entrypoint: python3 /main.py\n")
    file.write("    environment:\n")
    file.write("      - PYTHONUNBUFFERED=1\n")
    file.write("      - LOGGING_LEVEL=DEBUG\n")
    file.write("    networks:\n")
    file.write("      - testing_net\n")


def client(file, id):
    """
    Writes the client section of the docker-compose file for the given id
    """
    file.write(f"  client{id}:\n")
    file.write(f"    container_name: client{id}\n")
    file.write("    image: client:latest\n")
    file.write("    entrypoint: /client\n")
    file.write("    environment:\n")
    file.write(f"      - CLI_ID={id}\n")
    file.write("      - CLI_LOG_LEVEL=DEBUG\n")
    file.write("    networks:\n")
    file.write("      - testing_net\n")
    file.write("    depends_on:\n")
    file.write("      - server\n")


def networks(file):
    """
    Writes the network section of the docker-compose file
    """
    file.write("networks:\n")
    file.write("  testing_net:\n")
    file.write("    ipam:\n")
    file.write("      driver: default\n")
    file.write("      config:\n")
    file.write("        - subnet: 172.25.125.0/24\n")


if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python3 generate_dc.py n")
        print("Where n is the number of clients to generate")
        sys.exit(1)
    n = sys.argv[1]
    main(n)
