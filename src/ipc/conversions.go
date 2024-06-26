package ipc

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/setavenger/blindbitd/pb"
	"github.com/setavenger/blindbitd/src"
	"github.com/setavenger/blindbitd/src/logging"
	"github.com/setavenger/blindbitd/src/utils"
	"github.com/setavenger/go-bip352"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertWalletUTXOs(utxos []*src.OwnedUTXO, mapping src.LabelsMapping) []*pb.OwnedUTXO {
	var result []*pb.OwnedUTXO

	for _, utxo := range utxos {
		var label *pb.Label
		if utxo.Label != nil {
			label = convertLabel(utxo.Label, mapping)
		}
		result = append(result, &pb.OwnedUTXO{
			Txid:               utils.CopyBytes(utxo.Txid[:]),
			Vout:               utxo.Vout,
			Amount:             utxo.Amount,
			PubKey:             utils.CopyBytes(utxo.PubKey[:]),
			TimestampConfirmed: &timestamppb.Timestamp{Seconds: int64(utxo.Timestamp)},
			UtxoState:          convertState(utxo.State),
			Label:              label,
		})
	}

	return result
}

func convertLabel(label *bip352.Label, mapping src.LabelsMapping) *pb.Label {
	if label == nil {
		logging.WarningLogger.Println("label was nil")
		return nil
	}
	fullLabel := mapping.GetLabelByPubKey(label.PubKey)
	var comment string
	if fullLabel != nil {
		comment = fullLabel.Comment
	} else {
		logging.WarningLogger.Printf("There was no comment for label - m:%d pubkey: %x\n", label.M, label.PubKey)
	}

	var result = pb.Label{
		Address: label.Address,
		M:       label.M,
		Comment: comment,
	}

	return &result
}

func convertState(state src.UTXOState) pb.UTXOState {
	switch state {
	case src.StateUnconfirmed:
		return pb.UTXOState_UNCONFIRMED
	case src.StateUnspent:
		return pb.UTXOState_UNSPENT
	case src.StateSpent:
		return pb.UTXOState_SPENT
	case src.StateUnconfirmedSpent:
		return pb.UTXOState_SPENT_UNCONFIRMED
	default:
		return pb.UTXOState_UNKNOWN
	}
}

func convertToRecipients(recipients []*pb.TransactionRecipient) []*src.Recipient {
	var convertedRecipients []*src.Recipient

	for _, recipient := range recipients {
		convertedRecipients = append(convertedRecipients, &src.Recipient{
			Address:    recipient.Address,
			Amount:     int64(recipient.Amount),
			Annotation: recipient.Annotation,
		})
	}

	return convertedRecipients
}

func convertChainParam(params *chaincfg.Params) *pb.Chain {
	var chain pb.Chain

	switch params.Name {
	case chaincfg.MainNetParams.Name:
		chain.Chain = pb.ChainEnum_Mainnet
	case chaincfg.TestNet3Params.Name:
		chain.Chain = pb.ChainEnum_Testnet
	case chaincfg.SigNetParams.Name:
		chain.Chain = pb.ChainEnum_Signet
	case chaincfg.RegressionNetParams.Name:
		chain.Chain = pb.ChainEnum_Regtest
	}

	return &chain
}
