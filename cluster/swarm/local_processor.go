package swarm

import (
	"errors"
	"math/rand"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/docker/swarm/cluster"
	"github.com/docker/swarm/common"
)

// LocalProcessor is a synchronization task processor
type LocalProcessor struct {
	Cluster *Cluster
}

// Do task
func (t *LocalProcessor) Do(task *Task) (c string, err error) {
	logrus.Debugf("Local processor do task, task type:%d,container name:%s", task.TaskType, task.Container.Names)

	//containers create remove start stop
	//TODO if seize resource, needs stop container first, using multiple gorutine
	switch task.TaskType {
	case common.TaskTypeCreateContainer:
		c, err = t.createContainer(task.Container)
	case common.TaskTypeRemoveContainer:
		c, err = t.removeContainer(task.Container)
	case common.TaskTypeStartContainer:
		c, err = t.startContainer(task.Container)
	case common.TaskTypeStopContainer:
		c, err = t.stopContainer(task.Container)
	default:
		c = ""
		err = errors.New("unknow task type")
	}
	return c, err
}

func generateName(name string) string {
	return name + strconv.Itoa(rand.Int())
}

func (t *LocalProcessor) createContainer(container *cluster.Container) (c string, err error) {
	var newContainer *cluster.Container
	newContainer, err = t.Cluster.CreateContainer(container.Config, generateName(container.Names[0]), nil)
	if err != nil {
		logrus.Warnf("Scale up create container failed:%s", container.Names[0])
		return "", err
	}

	if err = t.Cluster.StartContainer(newContainer, nil); err != nil {
		logrus.Warnf("Scale up start container failed:%s", container.Names[0])
		return "", err
	}
	return newContainer.Names[0], nil
}

func (t *LocalProcessor) removeContainer(container *cluster.Container) (c string, err error) {
	//may be stop container first, this is force to remove container
	//remove volume or not remove volue, this method not remove volume
	if err = t.Cluster.RemoveContainer(container, true, false); err != nil {
		logrus.Warnf("remove container failed:%s", container.Names)
		return "", err
	}
	return container.Names[0], nil
}

func (t *LocalProcessor) startContainer(container *cluster.Container) (c string, err error) {
	return "", nil
}

func (t *LocalProcessor) stopContainer(container *cluster.Container) (c string, err error) {
	return "", nil
}
