"""global variables for all environments"""

load("group_vars/clusters.star", "clusters")

# the entrypoin
def module():
    return {
        "clusters": clusters 
    }
