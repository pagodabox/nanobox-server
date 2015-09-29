// Copyright (c) 2014 Pagoda Box Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public License,
// v. 2.0. If a copy of the MPL was not distributed with this file, You can
// obtain one at http://mozilla.org/MPL/2.0/.

//
package jobs

import (
	"github.com/nanobox-io/nanobox-boxfile"
	"github.com/nanobox-io/nanobox-golang-stylish"
	"github.com/nanobox-io/nanobox-server/config"
	"github.com/nanobox-io/nanobox-server/util"
)

//
type Restart struct {
	UID     string
	Success bool
}

// Proccess syncronies your docker containers with the boxfile specification
func (j *Restart) Process() {
	// add a lock so the service wont go down whil im running
	util.Lock()
	defer util.Unlock()
	
	j.Success = false

	util.LogInfo(stylish.Bullet("Restarting app in %s container...", j.UID))
	box := boxfile.NewFromPath("/vagrant/code/" + config.App + "/Boxfile")
	// restart payload
	payload := map[string]interface{}{
		"platform":    "local",
		"boxfile":     box.Node(j.UID).Parsed,
		"logtap_host": config.LogtapHost,
		"uid":         j.UID,
	}

	// run restart hook (blocking)
	if _, err := util.ExecHook("default-restart", j.UID, payload); err != nil {
		util.LogInfo("ERROR %v\n", err)
		return
	}

	j.Success = true
}
