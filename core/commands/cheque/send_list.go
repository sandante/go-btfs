package cheque

import (
	"fmt"
	"io"

	cmds "github.com/TRON-US/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
)

var ListSendChequesCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "List cheque(s) send to peers.",
	},

	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {

		listRet := ListChequeRet{}
		listRet.Cheques = make([]cheque, 0, 0)
		cheques, err := chain.SettleObject.SwapService.LastSendCheques()

		if err != nil {
			return err
		}
		for k, v := range cheques {
			var record cheque
			record.PeerID = k
			record.Beneficiary = v.Beneficiary.String()
			record.Vault = v.Vault.String()
			record.Payout = v.CumulativePayout

			listRet.Cheques = append(listRet.Cheques, record)
		}

		listRet.Len = len(listRet.Cheques)

		return cmds.EmitOnce(res, &listRet)
	},
	Type: ListChequeRet{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *ListChequeRet) error {
			fmt.Fprintf(w, "\t%-55s\t%-46s\t%-46s\tamount: \n", "peerID:", "vault:", "beneficiary:")
			for iter := 0; iter < out.Len; iter++ {
				fmt.Fprintf(w, "\t%-55s\t%-46s\t%-46s\t%d \n",
					out.Cheques[iter].PeerID,
					out.Cheques[iter].Vault,
					out.Cheques[iter].Beneficiary,
					out.Cheques[iter].Payout.Uint64(),
				)
			}

			return nil
		}),
	},
}
