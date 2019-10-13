# Memory

Memory is divided into three parts. 

They are ROM (Read Only Memory), RAM (Random Access Memory), and cache.

Comparsion of these memories:

\ | ROM | RAM | Cache
-- | --- | --- | -----
Speed | Fast | Fast | Very Fast
Read-Only | Y | N | N
Capacity | Tiny | Large | Small
Cost | Cheapest | Cheap | Expensive
Volatility | N | Y | Y

## ROM (Read Only Memory)

ROM is in short of *read only memory*. ROM's data is written when manufacturer make it. Generally, ROM stores important information of *system* or *hardware* permanently.

## Cache

In computing, a cache is a hardware or software component that stores data so that future requests for that data can be served faster; the data stored in a cache might be the result of an earlier computation or a copy of data stored elsewhere. 

Memory cache is a portion of the high-speed SRAM (static random access memory). With it, the computer avoids accessing low-speed DRAM (dynamic random access memory), making the computer perform faster and more efficiently. Today, most computers come with L3 cache or L2 cache, while older computers included only L1 cache.

## RAM (Random Access Memory)

**This part is mainly about DRAM (dynamic random access memory).**

RAM is the biggest part of memory, it is an extremely fast type of computer memory which temporarily stores all the information your PC needs right now and in the near future. 

### How RAM works?

**Memory is made up of bits (called *memory cell*) arranged in a two-dimensional grid.**

Memory cells are etched onto a silicon wafer in an array of columns (**bitlines**) and rows (**wordlines**). The intersection of a bitline and wordline constitutes the address of the memory cell.

DRAM works by sending a charge through the appropriate column (CAS) to activate the transistor at each bit in the column. When writing, the row lines contain the state the capacitor should take on. When reading, the sense-amplifier determines the level of charge in the capacitor. If it is more than 50 percent, it reads it as a 1; otherwise it reads it as a 0. The counter tracks the refresh sequence based on which rows have been accessed in what order. The length of time necessary to do all this is so short that it is expressed in **nanoseconds** (billionths of a second). A memory chip rating of 70ns means that it takes 70 nanoseconds to completely read and recharge each cell.

Memory cells alone would be worthless without some way to get information in and out of them. So the memory cells have a whole support infrastructure of other specialized circuits. These circuits perform functions such as:

- Identifying each row and column (row address select and column address select)
- Keeping track of the refresh sequence (counter)
- Reading and restoring the signal from a cell (sense amplifier)
- Telling a cell whether it should take a charge or not (write enable)
- Other functions of the memory controller include a series of tasks that include identifying the type, speed and amount of memory and checking for errors.

### Period

In DRAM's development, there are several periods:

*SDRAM*, *DDR*, *DDR2*, *DDR3* and *DDR4*.

There are comparison of these RAMs:

Standard | Internal rate (MHz) | Bus clock (MHz) | Prefetch | Data rate (MT/s) | Transfer rate (GB/s) | Voltage (V)
-- | -- | -- | --| --| --| --
SDRAM | 100-166 | 100-166 | 1n | 100-166 | 0.8-1.3 | 3.3
DDR | 133-200 | 133-200 | 2n | 266-400 | 2.1-3.2 | 2.5/2.6
DDR2 | 133-200 | 266-400 | 4n | 533-800 | 4.2-6.4 |1.8
DDR3 | 133-200 | 533-800 | 8n | 1066-1600 | 8.5-14.9 | 1.35/1.5
DDR4 | 133-200 | 1066-1600 | 8n | 2133-3200 | 17-21.3 | 1.2

### SDRAM

SDRAM, in short of *synchronous dynamic random-access memory*, is any DRAM where the operation of its external pin interface is coordinated by an externally supplied clock signal.

SDRAM is designed to synchronize itself with the timing of the CPU, which **enables the memory controller to know the exact clock cycle when the requested data will be ready**, so the CPU no longer has to wait between memory accesses. For example, PC66 SDRAM runs at 66 MT/s, PC100 SDRAM runs at 100 MT/s, PC133 SDRAM runs at 133 MT/s, and so on.

For SDRAM, the I/O, internal clock and bus clock are the same, so it also stand for SDR SDRAM (Single Data Rate SDRAM). **Single Data Rate means that SDR SDRAM can only read/write one time in a clock cycle.** SDRAM have to wait for the completion of the previous command to be able to do another read/write operation.

### DDR (Double Data Rate SDRAM)

The next generation of SDRAM is DDR, which achieves greater bandwidth than the preceding single data rate SDRAM by **transferring data on the rising and falling edges of the clock signal (double pumped)**. Effectively, it doubles the transfer rate without increasing the frequency of the clock. 

The transfer rate of DDR SDRAM is the double of SDR SDRAM without changing the internal clock. DDR SDRAM, as the first generation of DDR memory, the ***prefetch*** buffer is 2bit, which is the double of SDR SDRAM. 

> - Prefetch
Transfer data from main memory to temporary storage in readiness for later use.

### DDR2 (Double Data Rate 2 SDRAM)

Its primary benefit is the ability to **operate the external data bus twice as fast as DDR SDRAM**. This is achieved by improved bus signal. The prefetch buffer of DDR2 is 4 bit(double of DDR SDRAM). DDR2 memory is at the same internal clock speed (133~200MHz) as DDR,  but the transfer rate of DDR2 can reach 533~800 MT/s with the improved I/O bus signal. DDR2 533 and DDR2 800 memory types are on the market.

### DDR3 (Double Data Rate 3 SDRAM)

DDR3 memory reduces 40% power consumption compared to current DDR2 modules, allowing for lower operating currents and voltages (1.5 V, compared to DDR2's 1.8 V or DDR's 2.5 V). The transfer rate of DDR3 is 800~1600 MT/s. DDR3's prefetch buffer width is 8 bit, whereas DDR2's is 4 bit, and DDR's is 2 bit. DDR3 also adds two functions, such as ASR (Automatic Self-Refresh) and SRT (Self-Refresh Temperature). They can make the memory control the refresh rate according to the temperature variation.

### DDR4 (Double Data Rate 4 SDRAM)

DDR4 SDRAM provides the lower operating voltage (1.2V) and higher transfer rate. The transfer rate of DDR4 is 2133~3200 MT/s. DDR4 adds four new Bank Groups technology. Each bank group has the feature of singlehanded operation. DDR4 can process 4 data within a clock cycle, so DDR4's efficiency is better than DDR3 obviously. DDR4 also adds some functions, such as DBI (Data Bus Inversion), CRC (Cyclic Redundancy Check) and CA parity. They can enhance DDR4 memory's signal integrity, and improve the stability of data transmission/access.