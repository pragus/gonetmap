package gonetmap

const netmapDev = "/dev/netmap"

const (
	NetmapOpenNoMmap  = 0x040000 /* reuse mmap from parent */
	NetmapOpenIfname  = 0x080000 /* nr_name, nr_ringid, nr_flags */
	NetmapOpenArg1    = 0x100000
	NetmapOpenArg2    = 0x200000
	NetmapOpenArg3    = 0x400000
	NetmapOpenRingCfg = 0x800000 /* tx|rx rings|slots */
)
const NR_REG_MASK = 0xf
const RingTimestamp = 0x0002 /* set timestamp on *sync() */
const RingForward = 0x0004   /* enable NS_FORWARD for ring */

const NetmapHwRing = 0x4000   /* single NIC ring pair */
const NetmapSwRing = 0x2000   /* only host ring pair */
const NetmapRingMask = 0x0fff /* the ring number */
const NetmapNoTxPoll = 0x1000 /* no automatic txsync on poll */
const NetmapDoRxPoll = 0x8000 /* DO automatic rxsync on poll */

const NetmapRingMonitorTx = 0x100
const NetmapRingMonitorRx = 0x200
const NetmapRingZcopyMon = 0x400
const NetmapRingExclusive = 0x800     /* request exclusive access to the selected rings */
const NetmapRingPtNetmapHost = 0x1000 /* request ptnetmap host support */
const NetmapRingRxRingsOnly = 0x2000
const NetmapRingTxRingsOnly = 0x4000
const NetmapRingAcceptVnetHdr = 0x8000

const (
	NetmapBdgAttach     = 1  /* attach the NIC */
	NetmapBdgDetach     = 2  /* detach the NIC */
	NetmapBdgRegops     = 3  /* register bridge callbacks */
	NetmapBdgList       = 4  /* get bridge's info */
	NetmapBdgVnetHdr    = 5  /* set the port virtio-net-hdr length */
	NetmapBdgNewIf      = 6  /* create a virtual port */
	NetmapBdgDelIf      = 7  /* destroy a virtual port */
	NetmapPtHostCreate  = 8  /* create ptnetmap kthreads */
	NetmapPtHostDelete  = 9  /* delete ptnetmap kthreads */
	NetmapBdgPollingOn  = 10 /* delete polling kthread */
	NetmapBdgPollingOff = 11 /* delete polling kthread */
	NetmapVnetHdrGet    = 12 /* get the port virtio-net-hdr length */
	NetmapPoolsInfoGet  = 13 /* get memory allocator pools info */
)

const NetmapBdgHost = 1 /* attach the host stack on ATTACH */

type Register int

const (
	ReqDefault     Register = iota /* backward compat, should not be used. */
	ReqAllNic               = iota /* NR_REG_ALL_NIC, (default) all	hardware ring pairs */
	ReqSoftware             = iota /* NR_REG_SW, the ``host rings'', connecting to the	host stack. */
	ReqNicSoftware          = iota /* NR_REG_NIC_SW, all hardware rings and the host rings */
	ReqOneNic               = iota /* NR_REG_ONE_NIC, only the i-th	hardware ring pair, where the number is	in nr_ringid*/
	ReqPipeMaster           = iota /* NR_REG_PIPE_MASTER */
	ReqPipeSlave            = iota /* NR_REG_PIPE_SLAVE */
)

type Direction int

const (
	RX Direction = iota
	TX
)
