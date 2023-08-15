package reputation

import (
	"fmt"

	"github.com/stnokott/sot-api/internal/api/structs"
)

type nodeImplReputations structs.Reputations

// Label always returns "root" since there should only be one instance of this struct in a tree.
func (n nodeImplReputations) Label() string {
	return "root"
}

func (n nodeImplReputations) Children() []nodeProvider {
	nodes := make([]nodeProvider, len(n))
	i := 0
	for k, v := range n {
		nodes[i] = nodeImplReputation{
			label:      repTypeDisplaynames[k],
			Reputation: v,
		}
		i++
	}
	return nodes
}

type nodeImplReputation struct {
	label string
	*structs.Reputation
}

func (n nodeImplReputation) Label() string {
	return n.label
}

func (n nodeImplReputation) Children() (items []nodeProvider) {
	if n.Emblems != nil {
		items = make([]nodeProvider, len(n.Emblems))
		for i, v := range n.Emblems {
			items[i] = nodeImplEmblem(v)
		}
	} else if n.Campaigns != nil {
		items = make([]nodeProvider, len(n.Campaigns))
		i := 0
		for _, v := range n.Campaigns {
			items[i] = nodeImplCampaign(v)
			i++
		}
	}
	return
}

type nodeImplCampaign structs.Campaign

func (n nodeImplCampaign) Label() string {
	return n.Title
}

func (n nodeImplCampaign) Children() []nodeProvider {
	children := make([]nodeProvider, len(n.Emblems))
	for i, v := range n.Emblems {
		children[i] = nodeImplEmblem(v)
	}
	return children
}

type nodeImplEmblem structs.Emblem

func (n nodeImplEmblem) Label() string {
	return fmt.Sprintf("%s (%s)", n.Title, n.Subtitle)
}

func (n nodeImplEmblem) Children() []nodeProvider {
	return nil
}
