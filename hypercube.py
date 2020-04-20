import numpy as np 
import sys

def hypercube(ndim, diagonal=0):
    """Recursively construct the edge-connectivity of a hypercube

    Parameters
    ----------
    ndim : int
        Dimension of the hypercube
    diagonal : bool
        Value of the diagonal
        If True, vertices are considered connected to themselves

    Returns
    -------
    ndarray, [2**ndim, 2**ndim], bool
        connectivity pattern of the hypercube
    """
    if ndim == 0:
        return np.array([[diagonal]])
    else:
        D = hypercube(ndim-1, diagonal)
        I = np.eye(len(D),  dtype=np.int32)
        return np.block([
            [D, I],
            [I, D],
        ])

def writeOutPut(ndim):
    cube = hypercube(ndim).tolist()
    return cube
    
# if __name__ == "__main__":
sys.stdout.write('%s\n'%(
        writeOutPut(
            int(sys.argv[1]))))