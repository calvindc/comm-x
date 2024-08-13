package models

import (
	"fmt"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/calvindc/comm-x/nodefound/config"
	"github.com/calvindc/comm-x/orm-db"
	"gorm.io/gorm"
	"math/big"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

var vdb *gorm.DB

func TestGetAccountFeePolicy(t *testing.T) {
	filepath.Walk(os.TempDir(), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasPrefix(filepath.Base(path), "nodefound") {
			return os.Remove(path)
		}
		return nil
	})

	dbPath := path.Join(os.TempDir(), fmt.Sprintf("nodefound%s.db", utils.RandomString(10)))
	t.Logf("dbpath=%s", dbPath)
	err := os.Remove(dbPath) //use one db
	if err != nil {
		t.Logf("remove err %s", err)
	}
	vdb, err := orm_db.SetupDB("sqlite3", dbPath, &AccountFee{}, &AccountTokenFee{}, &ChannelParticipantFee{}, &TokenFee{})
	if err != nil {
		t.Logf("setupDB err %s", err)
	}
	address := utils.NewRandomAddress()

	fee := GetAccountFeePolicy(address, vdb)
	if fee.FeePolicy != config.DefaultFeePolicy || fee.FeePercent != config.DefaultFeePercentPart {
		t.Errorf("not equal default")
		return
	}
	//give some-self define
	fee.FeePolicy = FeePolicyPercent
	fee.FeePercent = 0
	fee.FeeConstant = big.NewInt(30)

	err = UpdataAccountDefaultFeePolicy(address, fee, vdb)
	if err != nil {
		t.Errorf("UpdataAccountDefaultFeePolicy err %s", err)
		return
	}

	fee.FeePercent = 50
	err = UpdataAccountDefaultFeePolicy(address, fee, vdb)
	if err != nil {
		t.Errorf("UpdataAccountDefaultFeePolicy err %s", err)
		return
	}

	//============
	fee2 := GetAccountFeePolicy(address, vdb)
	if !reflect.DeepEqual(fee, fee2) {
		t.Error("not equal")
		return
	}

	//============
	token := utils.NewRandomAddress()
	err = UpdateAccountTokenFee(address, token, fee, vdb)
	if err != nil {
		t.Errorf("UpdateAccountTokenFee err %s", err)
		return
	}

	fee3, err := GetAccountTokenFee(address, token, vdb)
	if err != nil {
		t.Errorf("GetAccountTokenFee err= %s", err)
		return
	}

	if !reflect.DeepEqual(fee, fee3) {
		t.Errorf("equal for GetAccountTokenFee err")

		return
	}
	t.Log(fee3.FeePolicy, fee3.FeePercent, fee3.FeeConstant)
	//============
	err = DeleteAccountAllFeeRate(address, vdb)
	if err != nil {
		t.Errorf("DeleteAccountAllFeeRate err %s", err)
		return
	}

	fee = GetAccountFeePolicy(address, vdb)
	if fee.FeePolicy != config.DefaultFeePolicy || fee.FeePercent != config.DefaultFeePercentPart {
		t.Error("not equal default")
		return
	}

}
