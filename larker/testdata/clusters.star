"""clusters"""

# functions
def cluster(name, shards):
    return { 
        "name": name, 
        "shards": shards,
    }

def get_base_config(cluster):
    return {
        "persistence": get_persistence(cluster["shards"]),
        "stats": get_stats(cluster["name"]),
        "dynamicconfig": get_dynamic_config(cluster["name"]),
    }

def get_dynamic_config(name):
    config = {
        "client": "dynamic-configerator",
        "dynamic-configerator": {
            "namespaces": name,
        },
        "applicationidentifier": "server",
        "cachedir": "/var/cache/dynamic-configerator-config",
        "iswatchfileenabled": "true",
    }

    return config

def get_persistence(shards):
    persistence = {
        "numHistoryShards": shards,
        "defaultStore": "caas-default",
        "visibilityStore": "caas-visibility",
    }

    return persistence

def get_stats(service):
    stats = {
        "exportInterval": "500ms",
        "exporter": {
            "m3": {
                "env": ENVIRONMENT,
                "hostPort": "127.0.0.1:9052",
                "service": service,
            },
        },
    }
    return stats

def create_cluster_configuration():
    out = []
    for cluster in _CLUSTERS:
        out.append({cluster["name"]: get_base_config(cluster)})
    return out

# Data
_CLUSTERS = [
    cluster("prod-1", 8192),
    cluster("prod-2", 16384)
]

clusters = create_cluster_configuration()