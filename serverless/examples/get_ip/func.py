import psutil
import socket


def list_all_network_interfaces_with_ipv4():
    interfaces = psutil.net_if_addrs()
    result = []
    for interface, addrs in interfaces.items():
        ipv4 = None
        for addr in addrs:
            if addr.family == socket.AF_INET:
                ipv4 = addr.address
        result.append((interface, ipv4))
    return result


def main(params):
    interfaces_with_ipv4 = list_all_network_interfaces_with_ipv4()
    result = {
        "interfaces": interfaces_with_ipv4
    }
    return result