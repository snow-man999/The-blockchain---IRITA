package clitest

import (
	"fmt"
	"testing"

	clientkeys "github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/tests"
	"github.com/stretchr/testify/require"

	"github.com/bianjieai/irita/app"
	"github.com/bianjieai/irita/modules/guardian"
)

func TestIritaCLIAddProfiler(t *testing.T) {
	t.Parallel()
	f := InitFixtures(t)

	fooAddr := f.KeyAddress(keyFoo)
	barAddr := f.KeyAddress(keyBar)

	// start irita server
	proc := f.GDStart()
	defer proc.Stop(false)

	description := "test"

	success, _, stderr := f.TxAddProfiler(fooAddr.String(), barAddr.String(), description, "-y")
	require.True(f.T, success)
	require.Empty(f.T, stderr)

	tests.WaitForNextNBlocksTM(1, f.Port)
	// Ensure transaction tags can be queried
	searchResult := f.QueryTxs(1, 50, "message.action=add_profiler", fmt.Sprintf("message.sender=%s", fooAddr))
	require.Len(t, searchResult.Txs, 1)

	expGuardian := guardian.NewGuardian(description, guardian.Ordinary, barAddr, fooAddr)

	res := f.QueryProfilers()
	require.NotEmpty(f.T, res)
	require.Contains(f.T, res, expGuardian)

	success, _, stderr = f.TxDeleteProfiler(fooAddr.String(), barAddr.String(), "-y")
	require.True(f.T, success)
	require.Empty(f.T, stderr)

	tests.WaitForNextNBlocksTM(1, f.Port)
	// Ensure transaction tags can be queried
	searchResult = f.QueryTxs(1, 50, "message.action=delete_profiler", fmt.Sprintf("message.sender=%s", fooAddr))
	require.Len(t, searchResult.Txs, 1)

	res = f.QueryProfilers()
	require.NotEmpty(f.T, res)
	require.NotContains(f.T, res, expGuardian)

	// Cleanup testing directories
	f.Cleanup()
}

func TestIritaCLIAddTrustee(t *testing.T) {
	t.Parallel()
	f := InitFixtures(t)

	// start irita server
	proc := f.GDStart()
	defer proc.Stop(false)

	fooAddr := f.KeyAddress(keyFoo)
	barAddr := f.KeyAddress(keyBar)

	description := "test"

	success, _, stderr := f.TxAddTrustee(fooAddr.String(), barAddr.String(), description, "-y")
	require.True(f.T, success)
	require.Empty(f.T, stderr)

	expGuardian := guardian.NewGuardian(description, guardian.Ordinary, barAddr, fooAddr)

	tests.WaitForNextNBlocksTM(1, f.Port)
	// Ensure transaction tags can be queried
	searchResult := f.QueryTxs(1, 50, "message.action=add_trustee", fmt.Sprintf("message.sender=%s", fooAddr))
	require.Len(t, searchResult.Txs, 1)

	res := f.QueryTrustees()
	require.NotEmpty(f.T, res)
	require.Contains(f.T, res, expGuardian)

	success, _, stderr = f.TxDeleteTrustee(fooAddr.String(), barAddr.String(), "-y")
	require.True(f.T, success)
	require.Empty(f.T, stderr)

	tests.WaitForNextNBlocksTM(1, f.Port)
	// Ensure transaction tags can be queried
	searchResult = f.QueryTxs(1, 50, "message.action=delete_trustee", fmt.Sprintf("message.sender=%s", fooAddr))
	require.Len(t, searchResult.Txs, 1)

	res = f.QueryTrustees()
	require.NotEmpty(f.T, res)
	require.NotContains(f.T, res, expGuardian)

	// Cleanup testing directories
	f.Cleanup()
}

//___________________________________________________________________________________
// iritacli tx guardian

// TxAddProfiler is iritacli tx guardian add-profiler
func (f *Fixtures) TxAddProfiler(from, address, description string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx guardian add-profiler %v --keyring-backend=test --from=%s --address=%s --description=%s", f.IritaCLIBinary, f.Flags(), from, address, description)
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), clientkeys.DefaultKeyPass)
}

// TxAddTrustee is iritacli tx guardian add-trustee
func (f *Fixtures) TxAddTrustee(from, address, description string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx guardian add-trustee %v --keyring-backend=test --from=%s --address=%s --description=%s", f.IritaCLIBinary, f.Flags(), from, address, description)
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), clientkeys.DefaultKeyPass)
}

// TxDeleteProfiler is iritacli tx guardian delete-profiler
func (f *Fixtures) TxDeleteProfiler(from, address string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx guardian delete-profiler %v --keyring-backend=test --from=%s --address=%s", f.IritaCLIBinary, f.Flags(), from, address)
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), clientkeys.DefaultKeyPass)
}

// TxDeleteTrustee is iritacli tx guardian delete-trustee
func (f *Fixtures) TxDeleteTrustee(from, address string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx guardian  delete-trustee %v --keyring-backend=test --from=%s --address=%s", f.IritaCLIBinary, f.Flags(), from, address)
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), clientkeys.DefaultKeyPass)
}

// QueryProfiler is iritacli query guardian profilers
func (f *Fixtures) QueryProfilers() (result guardian.Profilers) {
	cmd := fmt.Sprintf("%s query guardian profilers --output=%s %v", f.IritaCLIBinary, "json", f.Flags())
	out, _ := tests.ExecuteT(f.T, cmd, "")
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &result)
	require.NoError(f.T, err, "out %v\n, err %v", out, err)
	return
}

// QueryTrustee is iritacli query guardian profilers
func (f *Fixtures) QueryTrustees() (result guardian.Trustees) {
	cmd := fmt.Sprintf("%s query guardian trustees --output=%s %v", f.IritaCLIBinary, "json", f.Flags())
	out, _ := tests.ExecuteT(f.T, cmd, "")
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &result)
	require.NoError(f.T, err, "out %v\n, err %v", out, err)
	return
}
