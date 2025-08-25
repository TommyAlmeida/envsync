package env

type DiffResult struct {
    Missing   []string        `json:"missing"`
    Extra     []string        `json:"extra"`
    Different map[string]Diff `json:"different"`
    Same      []string        `json:"same"`
}

type Diff struct {
    Source string `json:"source"`
    Target string `json:"target"`
}

func CompareEnvs(source, target EnvVars) DiffResult {
    result := DiffResult{
        Different: make(map[string]Diff, len(source)),
    }

    for key, sourceValue := range source {
        if targetValue, exists := target[key]; exists {
            if sourceValue != targetValue {
                result.Different[key] = Diff{
                    Source: sourceValue,
                    Target: targetValue,
                }
            } else {
                result.Same = append(result.Same, key)
            }
        } else {
            result.Missing = append(result.Missing, key)
        }
    }


    for key := range target {
        if _, exists := source[key]; !exists {
            result.Extra = append(result.Extra, key)
        }
    }

    return result
}