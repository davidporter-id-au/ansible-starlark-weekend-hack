zones = ["zone_1", "zone_2"]

def create_tasks(name):
    out = []
    for z in zones: 
        out.append({
            "name": name + z,
            "debug": {
                "msg": "task ... "
            }
        })
    return out

def module():
    return create_tasks("setup-kafka-stuff")
        
    
