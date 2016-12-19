// Provide UI for the whole tool
// There is a main window looking a bit like a dock
// with push buttons opening each window managing each
// function of the software

package main

import (
    "os"
    "time"
    "fmt"

   "github.com/therecipe/qt/core"
   "github.com/therecipe/qt/gui"
   "github.com/therecipe/qt/widgets"
    api "github.com/osrg/gobgp/api"
    "github.com/Matt-Texier/local-mitigation-agent/gobgpclient"
    "google.golang.org/grpc"
)

// data strcutures used by both API functions and UI
// BGP flowspec update structure as exported from UI

type BgpFsRule struct {
    DstPrefix string
    SrcPrefix string
    AddrFam string
    Port string
    SrcPort string
    DstPort string
    TcpFlags string
    IcmpType string
    IcmpCode string
    ProtoNumber string
    PacketLen string
    Dscp string
    IpFrag string
    Action string
}

var BgpFsActivLib = []BgpFsRule{
    {DstPrefix: "1.1.1.1/32", SrcPrefix: "2.2.2.2/32", AddrFam: "IPv4", Port: "8080",
     SrcPort: "80", DstPort: "443", TcpFlags: "syn", IcmpType: "", IcmpCode: "", ProtoNumber: "6",
     PacketLen: "1024", Dscp: "22", IpFrag: "", Action: "",},
    {DstPrefix: "3.3.3.3/32", SrcPrefix: "4.4.4.4/32", AddrFam: "IPv4", Port: "8080",
     SrcPort: "80", DstPort: "443", TcpFlags: "syn", IcmpType: "", IcmpCode: "", ProtoNumber: "6",
     PacketLen: "1024", Dscp: "22", IpFrag: "", Action: "",},
    {DstPrefix: "5.5.5.5/32", SrcPrefix: "6.6.6.6/32", AddrFam: "IPv4", Port: "8080",
     SrcPort: "80", DstPort: "443", TcpFlags: "syn", IcmpType: "", IcmpCode: "", ProtoNumber: "6",
     PacketLen: "1024", Dscp: "22", IpFrag: "", Action: "",},
}

var (
    editRuleSrcPrefixLineEdit *widgets.QLineEdit
    editRuleDstPrefixLineEdit *widgets.QLineEdit
    editRuleIcmpTypeLineEdit *widgets.QLineEdit
    editRuleIcmpCodeLineEdit *widgets.QLineEdit
    editRuleIpProtoLineEdit *widgets.QLineEdit
    editRulePortLineEdit *widgets.QLineEdit
    editRuleSrcPortLineEdit *widgets.QLineEdit
    editRuleDstPortLineEdit *widgets.QLineEdit
    editRuleTcpFlagFilterLine *widgets.QLineEdit
    editRuleLenLineEdit *widgets.QLineEdit
    editRuleDscpLineEdit *widgets.QLineEdit
    editRuleFragFilterLine *widgets.QLineEdit
    editRuleTree *widgets.QTreeWidget
    consoleWindow *widgets.QMainWindow
    flowspecWindow *widgets.QMainWindow
)


var client api.GobgpApiClient

var (
    windowFlowSpecCreated bool
    windowBgpConsoleCreated bool
)

func main() {
    // initialise boolean that tell us if sub-windows is already reated
    windowFlowSpecCreated = false
    windowBgpConsoleCreated = false

    // launch gobgp API client
    timeout := grpc.WithTimeout(time.Second)
    conn, rpcErr := grpc.Dial("localhost:50051", timeout, grpc.WithBlock(), grpc.WithInsecure())
    if rpcErr != nil {
        fmt.Printf("GoBGP is probably not running on the local server ... Please start gobgpd process !\n")
        fmt.Println(rpcErr)
        return
    }
    client = api.NewGobgpApiClient(conn)

    widgets.NewQApplication(len(os.Args), os.Args)
    var dockWindow = widgets.NewQMainWindow(nil, 0)
    dockWindow.Layout().DestroyQObject()
    dockWindow.SetGeometry(core.NewQRect4(100, 100, 400, 50))
    dockWindow.SetWindowTitle("Gabu")
    var dockMainLayout = widgets.NewQHBoxLayout()
    dockMainLayout.SetSpacing(6)
    dockMainLayout.SetContentsMargins(11, 11, 11, 11)
    dockWindow.SetLayout(dockMainLayout)
    // main window "dock" push button
    var dockConsolePush = widgets.NewQPushButton2("GoBgp Console", dockWindow)
    var dockFlowSpecPush = widgets.NewQPushButton2("FlowSpec RIB", dockWindow)

    var dockButtonSizePolicy = widgets.NewQSizePolicy()
    dockButtonSizePolicy.SetHorizontalPolicy(widgets.QSizePolicy__Expanding)
    dockButtonSizePolicy.SetVerticalPolicy(widgets.QSizePolicy__Expanding)
    dockButtonSizePolicy.SetHorizontalStretch(0)
    dockButtonSizePolicy.SetVerticalStretch(0)
    dockConsolePush.SetSizePolicy(dockButtonSizePolicy)
    dockFlowSpecPush.SetSizePolicy(dockButtonSizePolicy)

    // Connect buttons to functions
    dockConsolePush.ConnectClicked(func(_ bool) { dockConsolButtonClicked() })
    dockFlowSpecPush.ConnectClicked(func(_ bool) { dockFspecButtonPushed() })
    // add button to main layout
    dockMainLayout.AddWidget(dockConsolePush, 0, 0)
    dockMainLayout.AddWidget(dockFlowSpecPush, 0, 0)
    dockWindow.Show()
    widgets.QApplication_Exec()

}

func dockConsolButtonClicked() {
    if(windowBgpConsoleCreated) {
        consoleWindow.Raise()

    } else {
        consoleWin()
        windowBgpConsoleCreated = true
    }
}

func dockFspecButtonPushed() {
    if(windowFlowSpecCreated) {
        flowspecWindow.Raise()
    } else {
        flowspecWin()
        windowFlowSpecCreated = true
    }

}





func consoleWin() {

    consoleWindow = widgets.NewQMainWindow(nil, 0)
    consoleWindow.Layout().DestroyQObject()
    consoleWindow.SetGeometry(core.NewQRect4(100, 100, 1000, 600))
    consoleWindow.SetWindowTitle("GoBGP Console")
    var mainLayout = widgets.NewQHBoxLayout()
    mainLayout.SetSpacing(6)
    mainLayout.SetContentsMargins(11, 11, 11, 11)
    consoleWindow.SetLayout(mainLayout)

    // console window widgets
    // log Frame
    var logFrame = widgets.NewQFrame(consoleWindow, 0)
    logFrame.SetFrameShape(widgets.QFrame__Panel)
    logFrame.SetFrameShadow(widgets.QFrame__Raised)
    var frameSizePolicy = widgets.NewQSizePolicy()
    frameSizePolicy.SetHorizontalPolicy(widgets.QSizePolicy__Preferred)
    frameSizePolicy.SetVerticalPolicy(widgets.QSizePolicy__Preferred)
    frameSizePolicy.SetHorizontalStretch(0)
    frameSizePolicy.SetVerticalStretch(0)
    logFrame.SetSizePolicy(frameSizePolicy)

    // layout for log
    var logLayout = widgets.NewQVBoxLayout()
    logLayout.SetSpacing(6);


    // Console text edit / display
    var logLabel = widgets.NewQLabel2("Console output", logFrame, 0)
    var logLabelSizePolicy = widgets.NewQSizePolicy()
    logLabelSizePolicy.SetHorizontalPolicy(widgets.QSizePolicy__Preferred)
    logLabelSizePolicy.SetVerticalPolicy(widgets.QSizePolicy__Preferred)
    logLabelSizePolicy.SetHorizontalStretch(0)
    logLabelSizePolicy.SetVerticalStretch(0)
    logLabelSizePolicy.SetHeightForWidth(logLabel.HasHeightForWidth())
    logLabel.SetSizePolicy(logLabelSizePolicy)
    logLabel.SetAlignment(core.Qt__AlignLeading|core.Qt__AlignLeft|core.Qt__AlignVCenter)
    logLayout.AddWidget(logLabel, 0, core.Qt__AlignLeft)

    var logText = widgets.NewQTextEdit(logFrame)
    var fixeFont = gui.NewQFont2("monospace", 10, 0, false)
    logText.SetFont(fixeFont)
    var logTextSizePolicy = widgets.NewQSizePolicy()
    logTextSizePolicy.SetHorizontalPolicy(widgets.QSizePolicy__Expanding)
    logTextSizePolicy.SetVerticalPolicy(widgets.QSizePolicy__Expanding)
    logTextSizePolicy.SetHorizontalStretch(0)
    logTextSizePolicy.SetVerticalStretch(0)
    logTextSizePolicy.SetHeightForWidth(logText.HasHeightForWidth())
    logText.SetSizePolicy(logTextSizePolicy)
    logLayout.AddWidget(logText, 0, 0)

    logFrame.SetLayout(logLayout)
    mainLayout.AddWidget(logFrame, 0, 0)

    // command Frame
    var cmdFrame = widgets.NewQFrame(consoleWindow, 0)
    cmdFrame.SetFrameShape(widgets.QFrame__Panel)
    cmdFrame.SetFrameShadow(widgets.QFrame__Raised)
    cmdFrame.SetSizePolicy(frameSizePolicy)

    // push buttons
    var cmdButtonSizePolicy = widgets.NewQSizePolicy()
    cmdButtonSizePolicy.SetHorizontalPolicy(widgets.QSizePolicy__Minimum)
    cmdButtonSizePolicy.SetVerticalPolicy(widgets.QSizePolicy__Fixed)
    cmdButtonSizePolicy.SetHorizontalStretch(0)
    cmdButtonSizePolicy.SetVerticalStretch(0)

    var cmdLabel = widgets.NewQLabel2("Basic Commands", cmdFrame, 0)
    cmdButtonSizePolicy.SetHeightForWidth(cmdLabel.HasHeightForWidth())
    cmdLabel.SetSizePolicy(cmdButtonSizePolicy)

    var    cmdNeighButton = widgets.NewQPushButton2("Neighbors", cmdFrame)
    cmdButtonSizePolicy.SetHeightForWidth(cmdNeighButton.HasHeightForWidth())
    cmdNeighButton.SetSizePolicy(cmdButtonSizePolicy)

    var cmdFsrib4Button = widgets.NewQPushButton2("IPv4 FlowSpec RIB", cmdFrame)
    cmdButtonSizePolicy.SetHeightForWidth(cmdFsrib4Button.HasHeightForWidth())
    cmdFsrib4Button.SetSizePolicy(cmdButtonSizePolicy)

    var cmdFsrib6Button = widgets.NewQPushButton2("IPv6 FlowSpec RIB", cmdFrame)
    cmdButtonSizePolicy.SetHeightForWidth(cmdFsrib6Button.HasHeightForWidth())
    cmdFsrib6Button.SetSizePolicy(cmdButtonSizePolicy)

    // layout for buttons
    var cmdLayout = widgets.NewQVBoxLayout()
    cmdLayout.AddWidget(cmdLabel, 0, core.Qt__AlignCenter)
    cmdLayout.AddWidget(cmdNeighButton, 0, 0)
    cmdLayout.AddWidget(cmdFsrib4Button, 0, 0)
    cmdLayout.AddWidget(cmdFsrib6Button, 0, 0)
    var cmdVerticalSpacer = widgets.NewQSpacerItem(20, 40, widgets.QSizePolicy__Fixed, widgets.QSizePolicy__Expanding)
    cmdLayout.AddItem(cmdVerticalSpacer)
    cmdFrame.SetLayout(cmdLayout)
    mainLayout.AddWidget(cmdFrame, 0, 0)

    // Connect push buttons
    cmdNeighButton.ConnectClicked(func(_ bool) { cmdNeighButtonClicked(logText) })
    cmdFsrib4Button.ConnectClicked(func(_ bool) { cmdFsrib4ButtonClicked(logText) })
    cmdFsrib6Button.ConnectClicked(func(_ bool) { cmdFsrib6ButtonClicked(logText) })
    consoleWindow.ConnectCloseEvent(consoleWindowClosed)
    consoleWindow.Show()
}

func consoleWindowClosed(event *gui.QCloseEvent){
    windowBgpConsoleCreated = false
}

func cmdNeighButtonClicked(logTextWidget *widgets.QTextEdit) {
    dumpNeigh := gobgpclient.TxtdumpGetNeighbor(client)

    for _, p := range dumpNeigh {
        logTextWidget.InsertPlainText(p)
    }
    logTextWidget.InsertPlainText("\n")
}

func cmdFsrib4ButtonClicked(logTextWidget *widgets.QTextEdit) {
    logTextWidget.InsertPlainText("Button FlowSpec 4\n\n")
}

func cmdFsrib6ButtonClicked(logTextWidget *widgets.QTextEdit) {
    logTextWidget.Append("Button FlowSpec 6\n\n")
}


func flowspecWin() {
    // Expanding Size policy
    var expandingSizePolicy = widgets.NewQSizePolicy()
    expandingSizePolicy.SetHorizontalPolicy(widgets.QSizePolicy__Expanding)
    expandingSizePolicy.SetVerticalPolicy(widgets.QSizePolicy__Expanding)
    expandingSizePolicy.SetHorizontalStretch(0)
    expandingSizePolicy.SetVerticalStretch(0)

    // preferred size policy
    var preferredSizePolicy = widgets.NewQSizePolicy()
    preferredSizePolicy.SetHorizontalPolicy(widgets.QSizePolicy__Preferred)
    preferredSizePolicy.SetVerticalPolicy(widgets.QSizePolicy__Preferred)
    preferredSizePolicy.SetHorizontalStretch(0)
    preferredSizePolicy.SetVerticalStretch(0)

    // Flowspec main window
    flowspecWindow = widgets.NewQMainWindow(nil, 0)
//    flowspecWindow.Layout().DestroyQObject()
    var flowspecCentralWid = widgets.NewQWidget(nil, 0)
    flowspecWindow.SetGeometry(core.NewQRect4(100, 100, 1000, 800))
    flowspecWindow.SetWindowTitle("Flowspec Configuration")
    var flowspecWindowLayout = widgets.NewQVBoxLayout()
    flowspecWindowLayout.SetSpacing(6)
    flowspecWindowLayout.SetContentsMargins(11, 11, 11, 11)
    flowspecCentralWid.SetLayout(flowspecWindowLayout)

    // create one frame and a dock, frame to host flwospec rule config
    // and a dock to manage flowspec Rib towards GoBGP
    var editRuleFrame = widgets.NewQFrame(flowspecWindow, 0)

    editRuleFrame.SetSizePolicy(preferredSizePolicy)

    editRuleFrame.SetFrameShape(widgets.QFrame__Panel)
    editRuleFrame.SetFrameShadow(widgets.QFrame__Raised)
    flowspecWindowLayout.AddWidget(editRuleFrame, 0, 0)

    var editRuleFrameLayout = widgets.NewQHBoxLayout()
    editRuleFrame.SetLayout(editRuleFrameLayout)


    // Create content of editRuleFrame
    // Widget for Tree that displays library
    var editRuleLibWid = widgets.NewQWidget(editRuleFrame, 0)
    editRuleLibWid.SetSizePolicy(preferredSizePolicy)
    editRuleFrameLayout.AddWidget(editRuleLibWid, 0, 0)
    var editRuleLibWidLayout = widgets.NewQVBoxLayout()
    editRuleLibWid.SetLayout(editRuleLibWidLayout)
    var editRuleLabel = widgets.NewQLabel2("Rules Library", editRuleLibWid, 0)
    editRuleTree = widgets.NewQTreeWidget(editRuleLibWid)
    editRuleTree.SetSizePolicy(expandingSizePolicy)
    editRuleLibWidLayout.AddWidget(editRuleLabel, 0, 0)
    editRuleLibWidLayout.AddWidget(editRuleTree, 0, 0)
    editRuleTree.SetColumnCount(13)
    var editRuleTreeHeaderItem = editRuleTree.HeaderItem()
    libHeaderLabels := []string{"Dst Prefix", "Src Prefix", "Port", "Src Port", "Dst Port", "TCP flags",
"ICMP Type", "ICMP code", "Proto Number", "Packet Len", "DSCP", "IP Frag", "Action"}
    for i, myLabel := range libHeaderLabels {
        editRuleTreeHeaderItem.SetText(i, myLabel)
    }
    fullfilTreeWithRuleLib(editRuleTree, BgpFsActivLib)
    var editRuleLibButtonFrame = widgets.NewQFrame(editRuleLibWid, 0)
    editRuleLibButtonFrame.SetFrameShape(widgets.QFrame__Panel)
    editRuleLibButtonFrame.SetFrameShadow(widgets.QFrame__Raised)
    editRuleLibWidLayout.AddWidget(editRuleLibButtonFrame, 0, 0)
    var editRuleLibButtonFrameLayout = widgets.NewQGridLayout2()
    editRuleLibButtonFrame.SetLayout(editRuleLibButtonFrameLayout)
    var (
        editRuleLibSaveButton = widgets.NewQPushButton2("Save library", editRuleLibButtonFrame)
        editRuleLibLoadButton = widgets.NewQPushButton2("Load library", editRuleLibButtonFrame)
        editRuleLibPushRibButton = widgets.NewQPushButton2("Push rule to BGP Rib", editRuleLibButtonFrame)
        editRuleLibConfRibButton = widgets.NewQPushButton2("Configure Rib", editRuleLibButtonFrame)
    )
    editRuleLibButtonFrameLayout.AddWidget(editRuleLibLoadButton, 0, 0, 0)
    editRuleLibButtonFrameLayout.AddWidget(editRuleLibSaveButton, 0, 1, 0)
    editRuleLibButtonFrameLayout.AddWidget(editRuleLibPushRibButton, 0, 2, 0)
    editRuleLibButtonFrameLayout.AddWidget(editRuleLibConfRibButton, 0, 3, 0)

    // var editRuleLibWidSpacer = widgets.NewQSpacerItem(20, 40, widgets.QSizePolicy__Expanding, widgets.QSizePolicy__Fixed)
    // editRuleLibButtonFrameLayout.AddItem(editRuleLibWidSpacer)

    // Edit rule widget creation: it includes all required
    // UI Widget to edit a BGP flowspec rule
    var editRuleMainWid = widgets.NewQWidget(editRuleFrame, 0)
    editRuleMainWid.SetSizePolicy(preferredSizePolicy)
    editRuleFrameLayout.AddWidget(editRuleMainWid, 0, core.Qt__AlignLeft)
    var editRuleMainWidLayout = widgets.NewQVBoxLayout()
    editRuleMainWid.SetLayout(editRuleMainWidLayout)
    // Editing widets of Edit Match filter
    var editRuleMainWidLabel = widgets.NewQLabel2("Edit Flowspec Match Filter", editRuleMainWid, 0)
    editRuleMainWidLayout.AddWidget(editRuleMainWidLabel, 0, 0)

    // Line edit for source and dest prefix
    var editRulePrefixGroupBox = widgets.NewQGroupBox2("Address family and Prefix filters", editRuleMainWid)
    editRuleMainWidLayout.AddWidget(editRulePrefixGroupBox, 0, 0)
    var editRulePrefixLayout = widgets.NewQGridLayout2()
    editRulePrefixGroupBox.SetLayout(editRulePrefixLayout)
    var (
        editRuleSrcPrefixLabel = widgets.NewQLabel2("Source Prefix:", editRulePrefixGroupBox, 0)
        editRuleDstPrefixLabel = widgets.NewQLabel2("Destination Prefix:", editRulePrefixGroupBox, 0)
        editAddrFamIpv4 = widgets.NewQRadioButton2("Flowspec IPv4", editRulePrefixGroupBox)
        editAddrFamIpv6 = widgets.NewQRadioButton2("Flowspec IPv6", editRulePrefixGroupBox)
    )
    editRuleSrcPrefixLineEdit = widgets.NewQLineEdit(nil)
    editRuleDstPrefixLineEdit = widgets.NewQLineEdit(nil)
    editRuleSrcPrefixLineEdit.SetPlaceholderText("1.1.1.1/32")
    editRuleDstPrefixLineEdit.SetPlaceholderText("2.2.2.2/24")
    editAddrFamIpv4.SetChecked(true)
    editRulePrefixLayout.AddWidget(editRuleSrcPrefixLabel, 1, 0, 0)
    editRulePrefixLayout.AddWidget(editRuleSrcPrefixLineEdit, 1, 1, 0)
    editRulePrefixLayout.AddWidget(editAddrFamIpv4, 0, 2, 0)
    editRulePrefixLayout.AddWidget(editRuleDstPrefixLabel, 0, 0, 0)
    editRulePrefixLayout.AddWidget(editRuleDstPrefixLineEdit, 0, 1, 0)
    editRulePrefixLayout.AddWidget(editAddrFamIpv6, 1, 2, 0)
    // horizontal widget to group together ICMP and proto type
    var editRuleIcmpProtoWid = widgets.NewQWidget(editRuleMainWid, 0)
    editRuleMainWidLayout.AddWidget(editRuleIcmpProtoWid, 0, 0)
    var editRuleIcmpProtoWidLayout = widgets.NewQHBoxLayout()
    editRuleIcmpProtoWidLayout.SetContentsMargins(0, 7, 0, 7)
    editRuleIcmpProtoWid.SetLayout(editRuleIcmpProtoWidLayout)
    // line edit for ICMP type and code
    var editRuleIcmpGroupBox = widgets.NewQGroupBox2("ICMP filters", editRuleMainWid)
    editRuleIcmpProtoWidLayout.AddWidget(editRuleIcmpGroupBox, 0, 0)
    var editRuleIcmpLayout = widgets.NewQGridLayout2()
    editRuleIcmpGroupBox.SetLayout(editRuleIcmpLayout)
    var (
        editRuleIcmpTypeLabel = widgets.NewQLabel2("ICMP Type:", editRuleIcmpGroupBox, 0)
        editRuleIcmpCodeLabel = widgets.NewQLabel2("ICMP Code:", editRuleIcmpGroupBox, 0)
    )
    editRuleIcmpTypeLineEdit = widgets.NewQLineEdit(nil)
    editRuleIcmpCodeLineEdit = widgets.NewQLineEdit(nil)
    editRuleIcmpTypeLineEdit.SetPlaceholderText("'=0' '=8'")
    editRuleIcmpCodeLineEdit.SetPlaceholderText("'=0'")
    editRuleIcmpLayout.AddWidget(editRuleIcmpTypeLabel, 0, 0, 0)
    editRuleIcmpLayout.AddWidget(editRuleIcmpTypeLineEdit, 0, 1, 0)
    editRuleIcmpLayout.AddWidget(editRuleIcmpCodeLabel, 1, 0, 0)
    editRuleIcmpLayout.AddWidget(editRuleIcmpCodeLineEdit, 1, 1, 0)
    // Line edit for IP protocol (Next header)
    var editRuleIpProtoGroupBox = widgets.NewQGroupBox2("IP protocol or Next header", editRuleMainWid)
    editRuleIcmpProtoWidLayout.AddWidget(editRuleIpProtoGroupBox, 0, 0)
    var editRuleIpProtoLayout = widgets.NewQGridLayout2()
    editRuleIpProtoGroupBox.SetLayout(editRuleIpProtoLayout)
    var (
        editRuleIpProtoLabel = widgets.NewQLabel2("Protocol number:", editRuleIcmpGroupBox, 0)
    )
    editRuleIpProtoLineEdit = widgets.NewQLineEdit(nil)
    editRuleIpProtoLineEdit.SetPlaceholderText("'=6' '=17'")
    editRuleIpProtoLayout.AddWidget(editRuleIpProtoLabel, 0, 0, 0)
    editRuleIpProtoLayout.AddWidget(editRuleIpProtoLineEdit, 0, 1, 0)

    // line edit for TCP/UDP ports
    var editRulePortGroupBox = widgets.NewQGroupBox2("Port filters", editRuleMainWid)
    editRuleMainWidLayout.AddWidget(editRulePortGroupBox, 0, 0)
    var editRulePortLayout = widgets.NewQGridLayout2()
    editRulePortGroupBox.SetLayout(editRulePortLayout)
    var (
        editRulePortLabel = widgets.NewQLabel2("Port:", editRulePortGroupBox, 0)
        editRuleSrcPortLabel = widgets.NewQLabel2("Src Port:", editRulePortGroupBox, 0)
        editRuleDstPortLabel = widgets.NewQLabel2("Dst Port:", editRulePortGroupBox, 0)
    )
    editRulePortLineEdit = widgets.NewQLineEdit(nil)
    editRuleSrcPortLineEdit = widgets.NewQLineEdit(nil)
    editRuleDstPortLineEdit = widgets.NewQLineEdit(nil)
    editRulePortLineEdit.SetPlaceholderText("'=80' '>=8080&<=8888'")
    editRuleSrcPortLineEdit.SetPlaceholderText("'=443&=80'")
    editRuleDstPortLineEdit.SetPlaceholderText("'>=1024&<=49151'")
    editRulePortLayout.AddWidget(editRulePortLabel, 0, 0, 0)
    editRulePortLayout.AddWidget(editRulePortLineEdit, 0, 1, 0)
    editRulePortLayout.AddWidget(editRuleSrcPortLabel, 1, 0, 0)
    editRulePortLayout.AddWidget(editRuleSrcPortLineEdit, 1, 1, 0)
    editRulePortLayout.AddWidget(editRuleDstPortLabel, 2, 0, 0)
    editRulePortLayout.AddWidget(editRuleDstPortLineEdit, 2, 1, 0)
    // line edit for TCP flags
    var editRuleTcpFlagGroupBox = widgets.NewQGroupBox2("TCP flags filter", editRuleMainWid)
    editRuleMainWidLayout.AddWidget(editRuleTcpFlagGroupBox, 0, 0)
    var editRuleTcpFlagLayout = widgets.NewQGridLayout2()
    editRuleTcpFlagGroupBox.SetLayout(editRuleTcpFlagLayout)
    var (
        editRuleTcpSynFlagCheck = widgets.NewQCheckBox2("SYN", editRuleTcpFlagGroupBox)
        editRuleTcpAckFlagCheck = widgets.NewQCheckBox2("ACK", editRuleTcpFlagGroupBox)
        editRuleTcpRstFlagCheck = widgets.NewQCheckBox2("RST", editRuleTcpFlagGroupBox)
        editRuleTcpFinFlagCheck = widgets.NewQCheckBox2("FIN", editRuleTcpFlagGroupBox)
        editRuleTcpPshFlagCheck = widgets.NewQCheckBox2("PSH", editRuleTcpFlagGroupBox)
        editRuleTcpEceFlagCheck = widgets.NewQCheckBox2("ECE", editRuleTcpFlagGroupBox)
        editRuleTcpUrgFlagCheck = widgets.NewQCheckBox2("URG", editRuleTcpFlagGroupBox)
        editRuleTcpCwrFlagCheck = widgets.NewQCheckBox2("CWR", editRuleTcpFlagGroupBox)
        editRuleLineSeparator = widgets.NewQFrame(editRuleTcpFlagGroupBox, 0)
        editRuleTcpOpAndCheck = widgets.NewQCheckBox2("AND", editRuleTcpFlagGroupBox)
        editRuleTcpOpNotCheck = widgets.NewQCheckBox2("NOT", editRuleTcpFlagGroupBox)
        editRuleTcpOpMatchCheck = widgets.NewQCheckBox2("MATCH", editRuleTcpFlagGroupBox)
        editRuleTcpFlagFilterLabel = widgets.NewQLabel2("Filter:", editRuleTcpFlagGroupBox, 0)
        editRuleTcpFlagAddButton = widgets.NewQPushButton2("Add", editRuleTcpFlagGroupBox)
    )
    editRuleTcpFlagFilterLine = widgets.NewQLineEdit(nil)
    editRuleLineSeparator.SetFrameShape(widgets.QFrame__VLine)
    editRuleLineSeparator.SetFrameShadow(widgets.QFrame__Sunken)
    editRuleTcpFlagLayout.AddWidget(editRuleTcpSynFlagCheck, 0, 0, 0)
    editRuleTcpFlagLayout.AddWidget(editRuleTcpAckFlagCheck, 0, 1, 0)
    editRuleTcpFlagLayout.AddWidget(editRuleTcpRstFlagCheck, 0, 2, 0)
    editRuleTcpFlagLayout.AddWidget(editRuleTcpFinFlagCheck, 0, 3, 0)
    editRuleTcpFlagLayout.AddWidget(editRuleTcpFlagFilterLabel, 0, 5, 0)
    editRuleTcpFlagLayout.AddWidget3(editRuleTcpFlagFilterLine, 0, 6, 1, 3, 0)
    editRuleTcpFlagLayout.AddWidget(editRuleTcpPshFlagCheck, 1, 0, 0)
    editRuleTcpFlagLayout.AddWidget(editRuleTcpEceFlagCheck, 1, 1, 0)
    editRuleTcpFlagLayout.AddWidget(editRuleTcpUrgFlagCheck, 1, 2, 0)
    editRuleTcpFlagLayout.AddWidget(editRuleTcpCwrFlagCheck, 1, 3, 0)
    editRuleTcpFlagLayout.AddWidget3(editRuleLineSeparator, 0, 4, 2, 1, 0)
    editRuleTcpFlagLayout.AddWidget(editRuleTcpOpAndCheck, 1, 5, 0)
    editRuleTcpFlagLayout.AddWidget(editRuleTcpOpNotCheck, 1, 6, 0)
    editRuleTcpFlagLayout.AddWidget(editRuleTcpOpMatchCheck, 1, 7, 0)
    editRuleTcpFlagLayout.AddWidget(editRuleTcpFlagAddButton, 1, 8, 0)

    // Line edit for packet length and DSCP
    var editRuleLenDscpGroupBox = widgets.NewQGroupBox2("Packet Length and DSCP", editRuleMainWid)
    editRuleMainWidLayout.AddWidget(editRuleLenDscpGroupBox, 0, 0)
    var editRuleLenDscpLayout = widgets.NewQGridLayout2()
    editRuleLenDscpGroupBox.SetLayout(editRuleLenDscpLayout)
    var (
        editRuleLenLabel = widgets.NewQLabel2("Packet length:", editRuleLenDscpGroupBox, 0)
        editRuleDscpLabel = widgets.NewQLabel2("DiffServ Codepoints:", editRuleLenDscpGroupBox, 0)

    )
    editRuleLenLineEdit = widgets.NewQLineEdit(nil)
    editRuleDscpLineEdit = widgets.NewQLineEdit(nil)
    editRuleLenLineEdit.SetPlaceholderText("'>=64&<=1024'")
    editRuleDscpLineEdit.SetPlaceholderText("'=46'")
    editRuleLenDscpLayout.AddWidget(editRuleLenLabel, 0, 0, 0)
    editRuleLenDscpLayout.AddWidget(editRuleLenLineEdit, 0, 1, 0)
    editRuleLenDscpLayout.AddWidget(editRuleDscpLabel, 0, 2, 0)
    editRuleLenDscpLayout.AddWidget(editRuleDscpLineEdit, 0, 3, 0)

    // Line edit and checkbox for fragment filtering
    var editRuleFragGroupBox = widgets.NewQGroupBox2("IP Fragment", editRuleMainWid)
    editRuleMainWidLayout.AddWidget(editRuleFragGroupBox, 0, 0)
    var editRuleFragLayout = widgets.NewQGridLayout2()
    editRuleFragGroupBox.SetLayout(editRuleFragLayout)
    var (
        editRuleIsfFragCheck = widgets.NewQCheckBox2("IsF", editRuleFragGroupBox)
        editRuleFfFragCheck = widgets.NewQCheckBox2("FF", editRuleFragGroupBox)
        editRuleLfFragCheck = widgets.NewQCheckBox2("LF", editRuleFragGroupBox)
        editRuleDfFragCheck = widgets.NewQCheckBox2("DF", editRuleFragGroupBox)
        editRuleAndFragCheck = widgets.NewQCheckBox2("AND", editRuleFragGroupBox)
        editRuleNotFragCheck = widgets.NewQCheckBox2("NOT", editRuleFragGroupBox)
        editRuleMatchFragCheck = widgets.NewQCheckBox2("MATCH", editRuleFragGroupBox)
        editRuleLineFragSeparator = widgets.NewQFrame(editRuleFragGroupBox, 0)
        editRuleFragFilterLabel = widgets.NewQLabel2("Filter:", editRuleFragGroupBox, 0)
        editRuleAddFragButton = widgets.NewQPushButton2("Add", editRuleFragGroupBox)
    )
    editRuleFragFilterLine = widgets.NewQLineEdit(nil)
    editRuleLineFragSeparator.SetFrameShape(widgets.QFrame__VLine)
    editRuleLineFragSeparator.SetFrameShadow(widgets.QFrame__Sunken)
    editRuleFragLayout.AddWidget(editRuleIsfFragCheck, 0, 0, 0)
    editRuleFragLayout.AddWidget(editRuleFfFragCheck, 0, 1, 0)
    editRuleFragLayout.AddWidget(editRuleLfFragCheck, 1, 0, 0)
    editRuleFragLayout.AddWidget(editRuleDfFragCheck, 1, 1, 0)
    editRuleFragLayout.AddWidget3(editRuleLineFragSeparator, 0, 2, 2, 1, 0)
    editRuleFragLayout.AddWidget(editRuleAndFragCheck, 1, 3, 0)
    editRuleFragLayout.AddWidget(editRuleNotFragCheck, 1, 4, 0)
    editRuleFragLayout.AddWidget(editRuleMatchFragCheck, 1, 5, 0)
    editRuleFragLayout.AddWidget(editRuleAddFragButton, 1, 6, 0)
    editRuleFragLayout.AddWidget3(editRuleFragFilterLabel, 0, 3, 1, 1, 0)
    editRuleFragLayout.AddWidget3(editRuleFragFilterLine, 0, 4, 1, 3, 0)

    // Editing widets of Action applied to match traffic
    var editRuleMainWidLabelMatch = widgets.NewQLabel2("Edit Flowspec Action", editRuleMainWid, 0)
    editRuleMainWidLayout.AddWidget(editRuleMainWidLabelMatch, 0, 0)
    // Match group box widget
    var editRuleActionGroupBox = widgets.NewQGroupBox2("Action applied", editRuleMainWid)
    editRuleMainWidLayout.AddWidget(editRuleActionGroupBox, 0, 0)
    var editRuleActionLayout = widgets.NewQGridLayout2()
    editRuleActionGroupBox.SetLayout(editRuleActionLayout)
    var (
        editRuleActionCombo = widgets.NewQComboBox(nil)
        editRuleRouteTargetLine = widgets.NewQLineEdit(nil)
    )
    editRuleActionCombo.AddItems([]string{"Drop", "Shape", "Redirect", "Marking"})
    editRuleActionLayout.AddWidget(editRuleActionCombo, 0, 0, 0)
    editRuleActionLayout.AddWidget(editRuleRouteTargetLine, 1, 0, 0)

    // global apply button
    var editRuleGlobButtonFrame = widgets.NewQFrame(editRuleMainWid, 0)
    var editRuleGlobButtonlayout = widgets.NewQGridLayout2()
    editRuleGlobButtonFrame.SetLayout(editRuleGlobButtonlayout)
    editRuleMainWidLayout.AddWidget(editRuleGlobButtonFrame, 0, 0)
    var (
        editGlobButtonNew = widgets.NewQPushButton2("New", editRuleGlobButtonFrame)
        editGlobButtonApply = widgets.NewQPushButton2("Apply", editRuleGlobButtonFrame)
        editGlobButtonReset = widgets.NewQPushButton2("Reset", editRuleGlobButtonFrame)
        editGlobButtonDelete = widgets.NewQPushButton2("Delete", editRuleGlobButtonFrame)

    )
    editRuleGlobButtonFrame.SetFrameShape(widgets.QFrame__Panel)
    editRuleGlobButtonFrame.SetFrameShadow(widgets.QFrame__Raised)
    editRuleGlobButtonlayout.AddWidget(editGlobButtonNew, 0, 0, 0)
    editRuleGlobButtonlayout.AddWidget(editGlobButtonApply, 0, 1, 0)
    editRuleGlobButtonlayout.AddWidget(editGlobButtonReset, 0, 2, 0)
    editRuleGlobButtonlayout.AddWidget(editGlobButtonDelete, 0, 3, 0)

    // var editRuleMainWidSpacer = widgets.NewQSpacerItem(20, 40, widgets.QSizePolicy__Fixed, widgets.QSizePolicy__Expanding)
    // editRuleMainWidLayout.AddItem(editRuleMainWidSpacer)
    // Connection of all widget to QT slots and actions
    // Tree Widget
    editRuleTree.ConnectItemClicked(editRuleLibItemSelected)
    // push button from rule edition
    // Connect push buttons
    editGlobButtonApply.ConnectClicked(func(_ bool) { editGlobButtonApplyFunc() })
    editGlobButtonNew.ConnectClicked(func(_ bool) { editGlobButtonNewFunc() })
    editGlobButtonDelete.ConnectClicked(func(_ bool) { editGlobButtonDeleteFunc() })
    editGlobButtonReset.ConnectClicked(func(_ bool) { editGlobButtonResetFunc() })

    // widget of the Rib tool dock
    var ribManipDock = widgets.NewQDockWidget("FlowSpec RIB tool", flowspecWindow, 0)
    // ribManipDock.SetSizePolicy(preferredSizePolicy)
    // flowspecWindowLayout.AddWidget(ribManipDock, 0, 0)
    flowspecWindow.AddDockWidget(core.Qt__BottomDockWidgetArea, ribManipDock)
    // main widget
    var ribManipDockWid = widgets.NewQWidget(nil, 0)
    var ribManipDockWidLayout = widgets.NewQHBoxLayout()
    ribManipDockWid.SetLayout(ribManipDockWidLayout)
    // Tree displaying BGP FS RIB
    var ribContentTree = widgets.NewQTreeWidget(ribManipDockWid)
    ribContentTree.SetSizePolicy(expandingSizePolicy)
    ribManipDockWidLayout.AddWidget(ribContentTree, 0, 0)
    // Buttons for rib manip
    var ribManipButtonWid = widgets.NewQWidget(ribManipDockWid, 0)
    var ribManipButtonWidLayout = widgets.NewQGridLayout2()
    ribManipButtonWid.SetLayout(ribManipButtonWidLayout)
    var (
        ribManipLoadButton = widgets.NewQPushButton2("Load/Reload BGP FS RIB", ribManipButtonWid)
        ribManipPushButton = widgets.NewQPushButton2("Push updates to RIB", ribManipButtonWid)
        ribManipAddRuleButton = widgets.NewQPushButton2("Add rule from Library", ribManipButtonWid)
        ribManipDeleteRuleButton = widgets.NewQPushButton2("Delete rule from RIB", ribManipButtonWid)
    )
    ribManipButtonWidLayout.AddWidget(ribManipLoadButton, 0, 0, 0)
    ribManipButtonWidLayout.AddWidget(ribManipPushButton, 1, 0, 0)
    ribManipButtonWidLayout.AddWidget(ribManipAddRuleButton, 2, 0, 0)
    ribManipButtonWidLayout.AddWidget(ribManipDeleteRuleButton, 3, 0, 0)
    ribManipDockWidLayout.AddWidget(ribManipButtonWid, 0, 0)

    ribManipDock.SetWidget(ribManipDockWid)

    flowspecWindow.SetCentralWidget(flowspecCentralWid)
    flowspecWindow.ConnectCloseEvent(flowspecWindowClosed)

    flowspecWindow.Show()
}

func flowspecWindowClosed(event *gui.QCloseEvent){
    windowFlowSpecCreated = false
}

// Copy the content of a flowspec rule structure into a TreeItem widget

func createFullfilItemWithRule(ty int, myTree *widgets.QTreeWidget, myRule BgpFsRule) {
    var myItem = widgets.NewQTreeWidgetItem3(myTree, ty)
    myItem.SetText(0, myRule.DstPrefix)
    myItem.SetText(1, myRule.SrcPrefix)
    myItem.SetText(2, myRule.Port)
    myItem.SetText(3, myRule.SrcPort)
    myItem.SetText(4, myRule.DstPort)
    myItem.SetText(5, myRule.TcpFlags)
    myItem.SetText(6, myRule.IcmpType)
    myItem.SetText(7, myRule.IcmpCode)
    myItem.SetText(8, myRule.ProtoNumber)
    myItem.SetText(9, myRule.PacketLen)
    myItem.SetText(10, myRule.Dscp)
    myItem.SetText(11, myRule.IpFrag)
    myItem.SetText(12, myRule.Action)
}

func fullfilItemWithRule(ty int, myItem *widgets.QTreeWidgetItem, myRule BgpFsRule) {
    myItem.SetText(0, myRule.DstPrefix)
    myItem.SetText(1, myRule.SrcPrefix)
    myItem.SetText(2, myRule.Port)
    myItem.SetText(3, myRule.SrcPort)
    myItem.SetText(4, myRule.DstPort)
    myItem.SetText(5, myRule.TcpFlags)
    myItem.SetText(6, myRule.IcmpType)
    myItem.SetText(7, myRule.IcmpCode)
    myItem.SetText(8, myRule.ProtoNumber)
    myItem.SetText(9, myRule.PacketLen)
    myItem.SetText(10, myRule.Dscp)
    myItem.SetText(11, myRule.IpFrag)
    myItem.SetText(12, myRule.Action)
}

func fullfilTreeWithRuleLib(myTree *widgets.QTreeWidget, myRuleLib []BgpFsRule) {
    for i, myRule := range myRuleLib {
        createFullfilItemWithRule(i, myTree, myRule)
    }
}

func fullfilLineEditWithBgpFs(myRule BgpFsRule) {
    editRuleSrcPrefixLineEdit.SetText(myRule.SrcPrefix)
    editRuleDstPrefixLineEdit.SetText(myRule.DstPrefix)
    editRuleIcmpTypeLineEdit.SetText(myRule.IcmpType)
    editRuleIcmpCodeLineEdit.SetText(myRule.IcmpCode)
    editRuleIpProtoLineEdit.SetText(myRule.ProtoNumber)
    editRulePortLineEdit.SetText(myRule.Port)
    editRuleSrcPortLineEdit.SetText(myRule.SrcPort)
    editRuleDstPortLineEdit.SetText(myRule.DstPort)
    editRuleTcpFlagFilterLine.SetText(myRule.TcpFlags)
    editRuleLenLineEdit.SetText(myRule.PacketLen)
    editRuleDscpLineEdit.SetText(myRule.Dscp)
    editRuleFragFilterLine.SetText(myRule.IpFrag)
}

func fullfilBgpFsWithLineEdit(myIndex int) {
    BgpFsActivLib[myIndex].SrcPrefix = editRuleSrcPrefixLineEdit.Text()
    BgpFsActivLib[myIndex].DstPrefix =  editRuleDstPrefixLineEdit.Text()
    BgpFsActivLib[myIndex].IcmpType = editRuleIcmpTypeLineEdit.Text()
    BgpFsActivLib[myIndex].IcmpCode = editRuleIcmpCodeLineEdit.Text()
    BgpFsActivLib[myIndex].ProtoNumber = editRuleIpProtoLineEdit.Text()
    BgpFsActivLib[myIndex].Port = editRulePortLineEdit.Text()
    BgpFsActivLib[myIndex].SrcPort = editRuleSrcPortLineEdit.Text()
    BgpFsActivLib[myIndex].DstPort = editRuleDstPortLineEdit.Text()
    BgpFsActivLib[myIndex].TcpFlags = editRuleTcpFlagFilterLine.Text()
    BgpFsActivLib[myIndex].PacketLen = editRuleLenLineEdit.Text()
    BgpFsActivLib[myIndex].Dscp = editRuleDscpLineEdit.Text()
    BgpFsActivLib[myIndex].IpFrag = editRuleFragFilterLine.Text()
}

// fucntion when an lib item is clicked

func editRuleLibItemSelected(myItem *widgets.QTreeWidgetItem, column int) {
    index := editRuleTree.IndexOfTopLevelItem(myItem)
    fullfilLineEditWithBgpFs(BgpFsActivLib[index])
}

// function to manage glob push button

func editGlobButtonNewFunc() {
    var myFsRule BgpFsRule
    myFsRule.DstPrefix = "New"
    BgpFsActivLib = append(BgpFsActivLib, myFsRule)
    createFullfilItemWithRule(len(BgpFsActivLib)-1, editRuleTree, BgpFsActivLib[len(BgpFsActivLib)-1])
}

func editGlobButtonApplyFunc() {
    var myItem *widgets.QTreeWidgetItem
    myItem = editRuleTree.CurrentItem()
    index := editRuleTree.IndexOfTopLevelItem(myItem)
    fullfilBgpFsWithLineEdit(index)
    fullfilItemWithRule(index, myItem, BgpFsActivLib[index])
}

func editGlobButtonDeleteFunc() {
    var myItem *widgets.QTreeWidgetItem
    myItem = editRuleTree.CurrentItem()
    index := editRuleTree.IndexOfTopLevelItem(myItem)
     if(index >= 0 && index < editRuleTree.TopLevelItemCount()) {
        myItem = editRuleTree.TakeTopLevelItem(index)
     }
    fmt.Printf("index: %d\n", index)
    BgpFsActivLib = append(BgpFsActivLib[:index], BgpFsActivLib[index+1:]...)
}

func editGlobButtonResetFunc() {
    fmt.Printf("Reset button\n")
}