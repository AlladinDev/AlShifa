package dependencyinjection

import (
	"fmt"
	"log"
	"reflect"
	"slices"

	appInterfaces "github.com/AlladinDev/AlShifa/internal/shared/interfaces"
)

type Dependency struct {
	DependencyAvailable bool
	DependencyName      string
}

type ServiceEntry struct {
	Value any
	Type  reflect.Type
}

type ModuleEntry struct {
	Name          string
	ConstructorFn func(diContainer appInterfaces.IDependencyInjection)
	Dependencies  []Dependency
}

type DI struct {
	ModulesWithDependencies []*ModuleEntry
	services                map[string]*ServiceEntry
	modules                 []*ModuleEntry
}

//helper functions ------

func (d *DI) GetServiceType(key string) reflect.Type {
	return d.services[key].Type
}

func (d *DI) PrintservicesInDI() {
	println("--------------------------------------------------------------")
	fmt.Println("Printing services in di container total : ", len(d.services))
	for key := range d.services {
		fmt.Println(key)
	}
	println("---------------------------------------------------------------")
}

func NewDI() appInterfaces.IDependencyInjection {
	return &DI{services: make(map[string]*ServiceEntry)}
}

func (d *DI) RegisterModule(moduleName string, constructorFn func(diContainer appInterfaces.IDependencyInjection), dependencies []string) {
	moduleDependencies := []Dependency{}

	for _, de := range dependencies {
		moduleDependencies = append(moduleDependencies, Dependency{
			DependencyName:      de,
			DependencyAvailable: false,
		})
	}

	d.modules = append(d.modules, &ModuleEntry{
		Name:          moduleName,
		ConstructorFn: constructorFn,
		Dependencies:  moduleDependencies,
	})
}

func (d *DI) AddService(key string, service any) {
	t := reflect.TypeOf(service)

	if t.Kind() != reflect.Pointer {
		log.Fatal("pass service as *Dependency, not a value")
	}

	// Detect **T
	if t.Elem().Kind() == reflect.Pointer {
		log.Fatalf("expected *Dependency but got **Dependency for key %s", key)
	}

	d.services[key] = &ServiceEntry{
		Value: service,
		Type:  t,
	}

	//now here check if the service with this key is needed by some module call function for its initialization
	for _, dependentModule := range d.ModulesWithDependencies {
		for index, dependency := range dependentModule.Dependencies {
			if dependency.DependencyName == key {
				dependentModule.Dependencies[index].DependencyAvailable = true
				d.InitializeModule(dependentModule)
			}
		}
	}
}

func (d *DI) InitializeModule(module *ModuleEntry) {

	for _, dependency := range module.Dependencies {
		if !dependency.DependencyAvailable {
			return
		}
	}

	//here remove this module from modulesWithDependencies array
	for index, moduleToRemove := range d.ModulesWithDependencies {
		if moduleToRemove.Name == module.Name {
			d.ModulesWithDependencies = slices.Delete(d.ModulesWithDependencies, index, index+1)
			break
		}
	}

	module.ConstructorFn(d)

}

func (d *DI) GetService(key string) (any, reflect.Type) {
	entry, ok := d.services[key]
	if !ok {
		return nil, nil
	}
	return entry.Value, entry.Type
}

func (d *DI) ServiceExists(key string) bool {
	_, exists := d.services[key]
	return exists
}

func (d *DI) ModuleExists(moduleName string) bool {
	for _, module := range d.modules {
		if module.Name == moduleName {
			return true
		}
	}
	return false
}

func (d *DI) PrintModulesRegistered() {
	println("---------------------------------------------")
	fmt.Printf("Printing Modules registered in dependency injection container total : %d ", len(d.modules))
	fmt.Println()
	for _, module := range d.modules {
		fmt.Println(module.Name)
	}
	println("----------------------------------------------")
}

func (d *DI) Instantiatemodules() {
	//first loop over the modules to check which modules can be safely instantiated
	for _, module := range d.modules {
		///here first we have to check that module constructor type should be a function
		if reflect.TypeOf(module.ConstructorFn).Kind() != reflect.Func {
			log.Fatalf("expected modules to be in dependency modules array but got %s", reflect.TypeOf(module.ConstructorFn).Kind())
		}

		//if module doesnt have any dependencies so just initialize it
		if len(module.Dependencies) == 0 {
			module.ConstructorFn(d)
			continue
		}

		moduleCanBeBuild := true
		moduleAddedToWaitingList := false

		//now constructor fn wants di we can ge it as we have it but first we have to check whether we have other dependencies which the module needs
		//means maybe module wants some instance through getService method so we have to check it
		for dependencyIndex, dependency := range module.Dependencies {
			if !d.ServiceExists(dependency.DependencyName) {
				println()
				moduleCanBeBuild = false
				//as this module needs something and we dont have it so append it to this slice
				//first check if this is already added to wait list no need to add again

				for _, waitingModule := range d.ModulesWithDependencies {
					if waitingModule.Name == module.Name {
						//it means module is already added to wait list no need to again add it new
						continue
					}
				}
				if !moduleAddedToWaitingList {
					d.ModulesWithDependencies = append(d.ModulesWithDependencies, module)
					moduleAddedToWaitingList = true
				}
				continue
			}
			//as this dependency is available so write its availability status as true
			module.Dependencies[dependencyIndex].DependencyAvailable = true
		}

		//it means if all its dependencies are available then just call it
		if moduleCanBeBuild {

			module.ConstructorFn(d)
		}

	}
	//at this point if we have still some moduleWithDependencies existing it means they need some instances which are not created so panic now
	if len(d.ModulesWithDependencies) > 0 {
		println("printing modules whose dependencies are  not fulfilled yet because their instances needed were not present in di")
		for _, module := range d.ModulesWithDependencies {
			println("-------------------------------------")
			count := 0
			for i := 0; i < len(module.Dependencies); i++ {
				if !module.Dependencies[i].DependencyAvailable {
					count++
				}
			}
			fmt.Printf("%s needs %d dependencies that are not  fulfilled yet", module.Name, count)
			println()
			for _, dependency := range module.Dependencies {
				if !dependency.DependencyAvailable {
					fmt.Printf("dependency : %s", dependency.DependencyName)
					fmt.Println()
				}
			}

			println("-------------------------------------")

		}
		log.Fatal("fullfill these dependencies to continue")
	}

	fmt.Print("Modules Build Successfully")
	d.PrintModulesRegistered()
	d.PrintservicesInDI()

}
