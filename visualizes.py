import matplotlib.pyplot as plt
import numpy as np
import pandas as pd 
import matplotlib.patches as mpatches


datasets =[ 
    "./datasheet/ring.csv",
    "./datasheet/circularladdar.csv",
    "./datasheet/hypercube.csv",
    "./datasheet/complete.csv"
]

demo = {
    "ring":{
        "label":"ring",
        "color":"red"
    },
    "circularladdar":{
        "label":"circularladdar",
        "color":"green"
    },
    "hypercube":{
        "label":"hypercube",
        "color":"blue"
    },
    "complete":{
        "label":"complete",
        "color":"orange"
    }
}

def preprocess(dataset):
    data = pd.read_csv(dataset)

    dimensions, msgs , numsNode , time , numsEdges , numsYoyo = [] , [] , [] , [], [], []

    def sumMsgs(row):
        msgs = [ row["#YoDown"] , row["#Yes"] , row["#YesPrune"], row["#No"]]
        return sum(map( int , msgs))
    cnt = 0
    for index, row in data.iterrows():
        if len(row.Description.split("_")) != 3:
            name, key = row.Description.split("dim")
            
            dimensions.append(int(key))
            msgs.append(sumMsgs(row))
            numsNode.append(row["#numsNode"])
            numsEdges.append(row['#nusmEdges'])

            msSplit , sSplit , mSplit = row["Duration"].split("ms") , row["Duration"].split("s") , row["Duration"].split("m")  
            if (len(msSplit) == 2):
                t = float(msSplit[0])
                time.append(t)
            elif (len(sSplit) == 2 and len(mSplit) != 2):
                t = float(sSplit[0])* 1000
                time.append(t)
            elif (len(msSplit) == 1 and len(mSplit) == 2):
                m = float(mSplit[0]) * 60
                s = (float(mSplit[1][:-1]) + m ) * 1000
                time.append(s)
            numsYoyo.append(cnt)
            cnt = 0
        else:
            cnt += 1
    return dimensions, msgs , numsNode , time , numsEdges , numsYoyo

for data in datasets:
    ds = data.split("/")[-1].split(".")[0]
    numsNode, msgs , dimensions , time , numsEdges , numsYoyo= preprocess(data)
    demo[ds]["numsNode"] = numsNode
    demo[ds]["msgSize"] = msgs 
    demo[ds]["time"] = time
    demo[ds]["numsEdge"] = numsEdges
    demo[ds]["numsYoyo"] = numsYoyo


def nodesVsMsgs():
    fig = plt.figure()
    fig.suptitle( "N vs Msg Size" , fontsize=20)
    legends = []
    plt.xlabel('nums Node')
    plt.ylabel('Total messages')

    for k , v in demo.items():
        plt.plot( 
            v["numsNode"],
            v["msgSize"],
            label=k,
            color=v['color'])
        legends.append(
            mpatches.Patch(color=v['color'],
                            label=v['label']))

    plt.legend(handles=legends)
    fig.savefig("./figures/N_vs_Msgs.png")
    plt.clf()

def nodesVsTimes():
    fig = plt.figure()
    fig.suptitle( "N vs Time" , fontsize=20)
    legends = []
    plt.xlabel('nums Node')
    plt.ylabel('time consume')

    for k , v in demo.items():
        plt.plot( 
            v["numsNode"],
            v["time"],
            label=k,
            color=v['color'])
        legends.append(
            mpatches.Patch(color=v['color'],
                            label=v['label']))

    plt.legend(handles=legends)
    fig.savefig("./figures/N_vs_Time.png")
    plt.clf()

def EdgesVsYoyos():
    fig = plt.figure()
    fig.suptitle( "numsYoyo vs EdgeSize" , fontsize=20)
    legends = []
    plt.xlabel('nums Edge')
    plt.ylabel('Yoyo stages')

    for k , v in demo.items():
        plt.plot( 
            v["numsEdge"],
            v["numsYoyo"],
            label=k,
            color=v['color'])
        legends.append(
            mpatches.Patch(color=v['color'],
                            label=v['label']))

    plt.legend(handles=legends)
    fig.savefig("./figures/NumsYoyo_vs_Edges.png")
    plt.clf()

def NodesVsYoyos():
    fig = plt.figure()
    fig.suptitle( "numsYoyo vs NodeSize" , fontsize=20)
    legends = []
    plt.xlabel('nums Node')
    plt.ylabel('Yoyo stages')

    for k , v in demo.items():
        plt.plot( 
            v["numsNode"],
            v["numsYoyo"],
            label=k,
            color=v['color'])
        legends.append(
            mpatches.Patch(color=v['color'],
                            label=v['label']))

    plt.legend(handles=legends)
    fig.savefig("./figures/NumsYoyo_vs_Nodes.png")
    plt.clf()

def EdgeVsMsgs():
    fig = plt.figure()
    fig.suptitle( "M vs Msg Size" , fontsize=20)
    legends = []
    plt.xlabel('nums Edge')
    plt.ylabel('Total messages')

    for k , v in demo.items():
        plt.plot( 
            v["numsEdge"],
            v["msgSize"],
            label=k,
            color=v['color'])
        legends.append(
            mpatches.Patch(color=v['color'],
                            label=v['label']))

    plt.legend(handles=legends)
    fig.savefig("./figures/M_vs_Msgs.png")
    plt.clf()

def EdgeVsTimes():
    fig = plt.figure()
    fig.suptitle( "M vs Time" , fontsize=20)
    legends = []
    plt.xlabel('nums Edge')
    plt.ylabel('time consume')

    for k , v in demo.items():
        plt.plot( 
            v["numsEdge"],
            v["time"],
            label=k,
            color=v['color'])
        legends.append(
            mpatches.Patch(color=v['color'],
                            label=v['label']))

    plt.legend(handles=legends)
    fig.savefig("./figures/M_vs_Time.png")
    plt.clf()

nodesVsTimes()
EdgeVsMsgs()
EdgeVsTimes()
EdgesVsYoyos()
nodesVsMsgs()
NodesVsYoyos()

# print(numsNode, msgs , dimensions, time)
# t = np.arange(0.0, 2.0, 0.01)
# s = 1 + np.sin(2*np.pi*t)
# plt.plot(t, s)

# plt.title('About as simple as it gets, folks')
# plt.show()
