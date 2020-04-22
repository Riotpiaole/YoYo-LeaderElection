# Assignment 3 Report 

## Introduction 

- The project compose creating generic concurrent implementation of Yoyo algorithm. Furthermore the project expands the limitation on applying various classical topology
  1. Circular Laddar Topology

  2. Ring Topology 

  3. HyperCube Topology 
       
  4. Complete Graph 


## Result 

1. Some of the figure can be used to display the result 
    - Message vs Nums of Edge 
      - ![](figures/M_vs_Msgs.png)

    - Message vs Nums of Node 
      - ![](figures/N_vs_Msgs.png)

    - Time vs Nums Node
      - ![](figures/N_vs_Time.png)

    - Time vs Nums edge
      - ![](figures/M_vs_Time.png)

    - NumsYoyo vs Nums edge
      - ![](figures/NumsYoyo_vs_Edges.png)
    
    - NumsYoyo vs Nums Node
      - ![](figures/NumsYoyo_vs_Nodes.png)

2. We can see the message complexity has heavier connection toward edge size. It make sense as the Yoyo complexity is $2m + 2mlog(s)$ in which it appears complete graph has the highest run time and message among all topology. 


