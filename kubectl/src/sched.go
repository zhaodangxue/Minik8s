package command

import (
	"minik8s/global"
	"minik8s/listwatch"

	"github.com/spf13/cobra"
)

var schedCommand = &cobra.Command{
	Use:   "sched",
	Short: "Manage the scheduler strategy",
	Long: `Manage the scheduler strategy. The scheduler strategy is used to determine how pods are scheduled onto nodes.
The scheduler strategy can be set to either "random" or "min-mem" or "min-cpu".`,
	Run:                        RunSched,
	SuggestionsMinimumDistance: 1,
	SuggestFor:                 []string{"schedule", "scheduling"},
	Example:                    "kubectl sched MininumCpuStrategy",
}

func RunSched(cmd *cobra.Command, args []string) {
	// Do Stuff Here
	strategy := args[0]
	listwatch.Publish(global.StrategyUpdateTopic(), strategy)
}
