package pdu

type CommandStatus uint32

const (
	ESME_ROK              CommandStatus = 0x00000000
	ESME_RINVMSGLEN       CommandStatus = 0x00000001
	ESME_RINVCMDLEN       CommandStatus = 0x00000002
	ESME_RINVCMDID        CommandStatus = 0x00000003
	ESME_RINVBNDSTS       CommandStatus = 0x00000004
	ESME_RALYBND          CommandStatus = 0x00000005
	ESME_RINVPRTFLG       CommandStatus = 0x00000006
	ESME_RINVREGDLVFLG    CommandStatus = 0x00000007
	ESME_RSYSERR          CommandStatus = 0x00000008
	ESME_RINVSRCADR       CommandStatus = 0x0000000A
	ESME_RINVDSTADR       CommandStatus = 0x0000000B
	ESME_RINVMSGID        CommandStatus = 0x0000000C
	ESME_RBINDFAIL        CommandStatus = 0x0000000D
	ESME_RINVPASWD        CommandStatus = 0x0000000E
	ESME_RINVSYSID        CommandStatus = 0x0000000F
	ESME_RCANCELFAIL      CommandStatus = 0x00000011
	ESME_RREPLACEFAIL     CommandStatus = 0x00000013
	ESME_RMSGQFUL         CommandStatus = 0x00000014
	ESME_RINVSERTYP       CommandStatus = 0x00000015
	ESME_RINVNUMDESTS     CommandStatus = 0x00000033
	ESME_RINVDLNAME       CommandStatus = 0x00000034
	ESME_RINVDESTFLAG     CommandStatus = 0x00000040
	ESME_RINVSUBREP       CommandStatus = 0x00000042
	ESME_RINVESMCLASS     CommandStatus = 0x00000043
	ESME_RCNTSUBDL        CommandStatus = 0x00000044
	ESME_RSUBMITFAIL      CommandStatus = 0x00000045
	ESME_RINVSRCTON       CommandStatus = 0x00000048
	ESME_RINVSRCNPI       CommandStatus = 0x00000049
	ESME_RINVDSTTON       CommandStatus = 0x00000050
	ESME_RINVDSTNPI       CommandStatus = 0x00000051
	ESME_RINVSYSTYP       CommandStatus = 0x00000053
	ESME_RINVREPFLAG      CommandStatus = 0x00000054
	ESME_RINVNUMMSGS      CommandStatus = 0x00000055
	ESME_RTHROTTLED       CommandStatus = 0x00000058
	ESME_RINVSCHED        CommandStatus = 0x00000061
	ESME_RINVEXPIRY       CommandStatus = 0x00000062
	ESME_RINVDFTMSGID     CommandStatus = 0x00000063
	ESME_RX_T_APPN        CommandStatus = 0x00000064
	ESME_RX_P_APPN        CommandStatus = 0x00000065
	ESME_RX_R_APPN        CommandStatus = 0x00000066
	ESME_RQUERYFAIL       CommandStatus = 0x00000067
	ESME_RINVOPTPARSTREAM CommandStatus = 0x000000C0
	ESME_ROPTPARNOTALLWD  CommandStatus = 0x000000C1
	ESME_RINVPARLEN       CommandStatus = 0x000000C2
	ESME_RMISSINGOPTPARAM CommandStatus = 0x000000C3
	ESME_RINVOPTPARAMVAL  CommandStatus = 0x000000C4
	ESME_RDELIVERYFAILURE CommandStatus = 0x000000FE
	ESME_RUNKNOWNERR      CommandStatus = 0x000000FF
)

var STATUS_DESCRIPTION = map[CommandStatus]string{
	ESME_ROK:              "No Error",
	ESME_RINVMSGLEN:       "Message Length is invalid",
	ESME_RINVCMDLEN:       "Command Length is invalid",
	ESME_RINVCMDID:        "Invalid Command ID",
	ESME_RINVBNDSTS:       "Incorrect BIND Status for given command",
	ESME_RALYBND:          "ESME Already in Bound State",
	ESME_RINVPRTFLG:       "Invalid Priority Flag",
	ESME_RINVREGDLVFLG:    "<Desc Not Set>",
	ESME_RSYSERR:          "System Error",
	ESME_RINVSRCADR:       "Invalid Source Address",
	ESME_RINVDSTADR:       "Invalid Destination Address",
	ESME_RINVMSGID:        "Invalid Message ID",
	ESME_RBINDFAIL:        "Bind Failed",
	ESME_RINVPASWD:        "Invalid Password",
	ESME_RINVSYSID:        "Invalid System ID",
	ESME_RCANCELFAIL:      "Cancel SM Failed",
	ESME_RREPLACEFAIL:     "Replace SM Failed",
	ESME_RMSGQFUL:         "Message Queue is full",
	ESME_RINVSERTYP:       "Invalid Service Type",
	ESME_RINVNUMDESTS:     "Invalid number of destinations",
	ESME_RINVDLNAME:       "Invalid Distribution List name",
	ESME_RINVDESTFLAG:     "Invalid Destination Flag (submit_multi)",
	ESME_RINVSUBREP:       "Invalid Submit With Replace request (replace_if_present_flag set)",
	ESME_RINVESMCLASS:     "Invalid esm_class field data",
	ESME_RCNTSUBDL:        "Cannot submit to Distribution List",
	ESME_RSUBMITFAIL:      "submit_sm or submit_multi failed",
	ESME_RINVSRCTON:       "Invalid Source address TON",
	ESME_RINVSRCNPI:       "Invalid Source address NPI",
	ESME_RINVDSTTON:       "Invalid Destination address TON",
	ESME_RINVDSTNPI:       "Invalid Destination address NPI",
	ESME_RINVSYSTYP:       "Invalid system_type field",
	ESME_RINVREPFLAG:      "Invalid replace_if_present flag",
	ESME_RINVNUMMSGS:      "Invalid number of messages",
	ESME_RTHROTTLED:       "Throttling error (ESME has exceeded allowed message limits)",
	ESME_RINVSCHED:        "Invalid Scheduled Delivery Time",
	ESME_RINVEXPIRY:       "Invalid message validity period (Expiry Time)",
	ESME_RINVDFTMSGID:     "Predefined Message is invalid or not found",
	ESME_RX_T_APPN:        "ESME received Temporary App Error Code",
	ESME_RX_P_APPN:        "ESME received Permanent App Error Code",
	ESME_RX_R_APPN:        "ESME received Reject Message Error Code",
	ESME_RQUERYFAIL:       "query_sm request failed",
	ESME_RINVOPTPARSTREAM: "Error in the optional part of the PDU body",
	ESME_ROPTPARNOTALLWD:  "Optional Parameter not allowed",
	ESME_RINVPARLEN:       "Invalid Parameter Length",
	ESME_RMISSINGOPTPARAM: "Expected Optional Parameter missing",
	ESME_RINVOPTPARAMVAL:  "Invalid Optional Parameter Value",
	ESME_RDELIVERYFAILURE: "Delivery Failure (used data_sm_resp)",
	ESME_RUNKNOWNERR:      "Unknown Error",
}
