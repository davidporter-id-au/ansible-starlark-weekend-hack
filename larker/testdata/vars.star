"""global variables for all environments"""

load("testdata/clusters.star", "clusters")

def module():
    return {
        "clusters": clusters 
    }
