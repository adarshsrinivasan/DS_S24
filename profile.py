"""
Cloudlab profile to setup. Each node runs on Ubuntu 22.04.

Instructions:
Create an experiment in CloudLab.
At least have 7 nodes in the topology for the experiment.
Node1 = Postgres Server
Node2 = MongoDB Server
Node3 = Seller Server
Node4 = Buyer Server
Node5 = Seller Client
Node6 = Buyer Client
Node7 = Transaction Server (from assignment 2)

Wait for the profile instance to start, then click on the node in the topology and choose the `shell` menu item.
"""

# Import the Portal object.
import geni.portal as portal
# Import the ProtoGENI library.
import geni.rspec.pg as pg

BASE_IP = "10.20.1"
BANDWIDTH = 10000000
# Create a portal context.
pc = portal.Context()

pc.defineParameter(
    "nodeCount", "Number of nodes in the experiment.", portal.ParameterType.INTEGER, 7,
    longDescription="Number of nodes in the topology. It is recommended to keep it 7")

pc.defineParameter(
    "assignmentNumber", "Assignment number to handle deployment accordingly", portal.ParameterType.INTEGER, 1,
    longDescription="Assignment number to handle deployment accordingly. Expected: [1 | 2]")

params = pc.bindParameters()
# Create a Request object to start building the RSpec.
request = pc.makeRequestRSpec()

nodes = []
lan = request.LAN()
lan.bandwidth = BANDWIDTH

for i in range(params.nodeCount):
    # Add a raw PC to the request.
    name = "node"+str(i+1)
    node = request.RawPC(name)

    interface = node.addInterface("if1")
    interface.addAddress(pg.IPv4Address("{}.{}".format(BASE_IP, 1 + len(nodes)), "255.255.255.0"))
    lan.addInterface(interface)

    nodes.append(node)

for i, node in enumerate(nodes):
    # Install and execute a script that is contained in the repository.
    node.addService(pg.Execute(shell="sh", command="sudo /local/repository/start.sh {} {} > /local/repository/setup-{}.log 2>&1".format(i, params.assignmentNumber, i)))

# Print the RSpec to the enclosing page.
pc.printRequestRSpec(request)