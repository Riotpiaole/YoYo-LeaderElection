import sys
import numpy as np 
import networkx as nx

def generateRingTopologyAdjacencyMatrix(nodeSize):
    graph = nx.generators.classic.circulant_graph(nodeSize,[1])
    matrxiGraph = nx.adjacency_matrix(graph)
    return matrxiGraph.toarray()

def generatlollipopGraph(M, N): 
    graph = nx.generators.classic.circular_ladder_graph(nodeSize)
    matrxiGraph = nx.adjacency_matrix(graph)
    return matrxiGraph.toarray()


def generateCircularLadder(nodeSize):
    graph = nx.generators.classic.circular_ladder_graph(nodeSize)
    matrxiGraph = nx.adjacency_matrix(graph)
    return matrxiGraph.toarray()

def writeOutPut(func, ndim):
    topology = func(ndim).tolist()
    return topology
    
if __name__ == "__main__":
    if  sys.argv[1] == "ring":
        func = generateRingTopologyAdjacencyMatrix
    elif sys.argv[1] == "circularladder":
        func = generateCircularLadder


    sys.stdout.write('%s\n'%(
            writeOutPut(
                func, int(sys.argv[2]))))