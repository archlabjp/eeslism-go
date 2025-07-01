package eeslism

import (
	"testing"
)

// Equipment Test Framework
// 各機器の基本的な動作テストを実装

// TestBOI tests Boiler equipment
func TestBOI(t *testing.T) {
	t.Run("BOI_BasicOperation", func(t *testing.T) {
		// ボイラーの基本動作テスト
		testBoilerBasicOperation(t)
	})
	
	t.Run("BOI_EfficiencyCalculation", func(t *testing.T) {
		// ボイラー効率計算テスト
		testBoilerEfficiency(t)
	})
}

// TestREFA tests Refrigeration equipment (Heat Pump)
func TestREFA(t *testing.T) {
	t.Run("REFA_BasicOperation", func(t *testing.T) {
		// ヒートポンプの基本動作テスト
		testHeatPumpBasicOperation(t)
	})
	
	t.Run("REFA_COPCalculation", func(t *testing.T) {
		// COP計算テスト
		testHeatPumpCOP(t)
	})
}

// TestCOL tests Solar Collector equipment
func TestCOL(t *testing.T) {
	t.Run("COL_BasicOperation", func(t *testing.T) {
		// 太陽熱集熱器の基本動作テスト
		testCollectorBasicOperation(t)
	})
	
	t.Run("COL_EfficiencyCalculation", func(t *testing.T) {
		// 集熱効率計算テスト
		testCollectorEfficiency(t)
	})
}

// TestSTANK tests Storage Tank equipment
func TestSTANK(t *testing.T) {
	t.Run("STANK_BasicOperation", func(t *testing.T) {
		// 蓄熱槽の基本動作テスト
		testStorageTankBasicOperation(t)
	})
	
	t.Run("STANK_TemperatureStratification", func(t *testing.T) {
		// 温度成層テスト
		testStorageTankStratification(t)
	})
}

// TestHEX tests Heat Exchanger equipment
func TestHEX(t *testing.T) {
	t.Run("HEX_BasicOperation", func(t *testing.T) {
		// 熱交換器の基本動作テスト
		testHeatExchangerBasicOperation(t)
	})
	
	t.Run("HEX_EffectivenessCalculation", func(t *testing.T) {
		// 熱交換効率計算テスト
		testHeatExchangerEffectiveness(t)
	})
}

// TestHCC tests Heating/Cooling Coil equipment
func TestHCC(t *testing.T) {
	t.Run("HCC_BasicOperation", func(t *testing.T) {
		// 冷暖房コイルの基本動作テスト
		testCoilBasicOperation(t)
	})
	
	t.Run("HCC_HeatTransferCalculation", func(t *testing.T) {
		// 熱伝達計算テスト
		testCoilHeatTransfer(t)
	})
}

// TestPIPE tests Pipe equipment
func TestPIPE(t *testing.T) {
	t.Run("PIPE_BasicOperation", func(t *testing.T) {
		// 配管の基本動作テスト
		testPipeBasicOperation(t)
	})
	
	t.Run("PIPE_HeatLossCalculation", func(t *testing.T) {
		// 配管熱損失計算テスト
		testPipeHeatLoss(t)
	})
}

// TestDUCT tests Duct equipment
func TestDUCT(t *testing.T) {
	t.Run("DUCT_BasicOperation", func(t *testing.T) {
		// ダクトの基本動作テスト
		testDuctBasicOperation(t)
	})
	
	t.Run("DUCT_PressureLossCalculation", func(t *testing.T) {
		// ダクト圧力損失計算テスト
		testDuctPressureLoss(t)
	})
}

// TestPUMP tests Pump equipment
func TestPUMP(t *testing.T) {
	t.Run("PUMP_BasicOperation", func(t *testing.T) {
		// ポンプの基本動作テスト
		testPumpBasicOperation(t)
	})
	
	t.Run("PUMP_PowerCalculation", func(t *testing.T) {
		// ポンプ動力計算テスト
		testPumpPower(t)
	})
}

// TestFAN tests Fan equipment
func TestFAN(t *testing.T) {
	t.Run("FAN_BasicOperation", func(t *testing.T) {
		// ファンの基本動作テスト
		testFanBasicOperation(t)
	})
	
	t.Run("FAN_PowerCalculation", func(t *testing.T) {
		// ファン動力計算テスト
		testFanPower(t)
	})
}

// TestVAV tests Variable Air Volume equipment
func TestVAV(t *testing.T) {
	t.Run("VAV_BasicOperation", func(t *testing.T) {
		// VAVの基本動作テスト
		testVAVBasicOperation(t)
	})
	
	t.Run("VAV_FlowControlCalculation", func(t *testing.T) {
		// 風量制御計算テスト
		testVAVFlowControl(t)
	})
}

// TestSTHEAT tests Storage Heater equipment
func TestSTHEAT(t *testing.T) {
	t.Run("STHEAT_BasicOperation", func(t *testing.T) {
		// 電気蓄熱暖房器の基本動作テスト
		testStorageHeaterBasicOperation(t)
	})
	
	t.Run("STHEAT_ThermalStorageCalculation", func(t *testing.T) {
		// 蓄熱計算テスト
		testStorageHeaterThermalStorage(t)
	})
}

// TestTHEX tests Total Heat Exchanger equipment
func TestTHEX(t *testing.T) {
	t.Run("THEX_BasicOperation", func(t *testing.T) {
		// 全熱交換器の基本動作テスト
		testTotalHeatExchangerBasicOperation(t)
	})
	
	t.Run("THEX_EfficiencyCalculation", func(t *testing.T) {
		// 全熱交換効率計算テスト
		testTotalHeatExchangerEfficiency(t)
	})
}

// TestPV tests Photovoltaic equipment
func TestPV(t *testing.T) {
	t.Run("PV_BasicOperation", func(t *testing.T) {
		// 太陽光発電の基本動作テスト
		testPVBasicOperation(t)
	})
	
	t.Run("PV_PowerGenerationCalculation", func(t *testing.T) {
		// 発電量計算テスト
		testPVPowerGeneration(t)
	})
	
	t.Run("PV_TemperatureCorrection", func(t *testing.T) {
		// 温度補正計算テスト
		testPVTemperatureCorrection(t)
	})
}

// TestOMVAV tests Outside Mount VAV equipment
func TestOMVAV(t *testing.T) {
	t.Run("OMVAV_BasicOperation", func(t *testing.T) {
		// 集熱屋根用VAVの基本動作テスト
		testOMVAVBasicOperation(t)
	})
	
	t.Run("OMVAV_SolarCollectionControl", func(t *testing.T) {
		// 集熱制御テスト
		testOMVAVSolarControl(t)
	})
}

// TestDESI tests Desiccant equipment
func TestDESI(t *testing.T) {
	t.Run("DESI_BasicOperation", func(t *testing.T) {
		// デシカント空調機の基本動作テスト
		testDesiccantBasicOperation(t)
	})
	
	t.Run("DESI_MoistureAdsorptionCalculation", func(t *testing.T) {
		// 吸湿計算テスト
		testDesiccantMoistureAdsorption(t)
	})
}

// TestEVAC tests Evaporative Cooler equipment
func TestEVAC(t *testing.T) {
	t.Run("EVAC_BasicOperation", func(t *testing.T) {
		// 気化冷却器の基本動作テスト
		testEvaporativeCoolerBasicOperation(t)
	})
	
	t.Run("EVAC_CoolingEfficiencyCalculation", func(t *testing.T) {
		// 冷却効率計算テスト
		testEvaporativeCoolerEfficiency(t)
	})
}

// 統合テスト：複数機器の連携動作テスト
func TestEquipmentIntegration(t *testing.T) {
	t.Run("Integration_SolarSystem", func(t *testing.T) {
		// 太陽熱システム（COL + STANK + PUMP）の統合テスト
		testSolarSystemIntegration(t)
	})
	
	t.Run("Integration_HVACSystem", func(t *testing.T) {
		// 空調システム（REFA + HCC + FAN + VAV）の統合テスト
		testHVACSystemIntegration(t)
	})
	
	t.Run("Integration_PVSystem", func(t *testing.T) {
		// 太陽光発電システムの統合テスト
		testPVSystemIntegration(t)
	})
}