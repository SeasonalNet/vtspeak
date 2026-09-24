# Reverse-engineering notebooks

Notebooks are optional lab companions for questions where inspecting captured
values or plots helps form and test a hypothesis. The Markdown findings remain
the reviewed source of truth; analysis code belongs in the existing standalone
scripts when it should be reusable outside Jupyter.

## Tree feature observations

[`tree-feature-observations.ipynb`](tree-feature-observations.ipynb) matches
the three checked-in Stage 6 runtime lookup captures against local `tree3`
resources and visualizes observed raw selector-value variation. It reuses
`compare_tree3_runtime.py` and `tree3.py`. It does not assign phonetic labels
or establish complete input coverage.

Run it from the repository root with a Python 3 Jupyter environment:

```sh
jupyter lab docs/reverse-engineering/notebooks/tree-feature-observations.ipynb
```

The notebook reads the tracked logs and local `data-common/dict-eng` and
`data-paul/M16/ttsdata/tree3` inputs. Keep those vendor inputs read-only. It
does not embed the inputs or save executed outputs in Git; execute it locally
to render the summary and SVG heatmap. No project or production dependencies
are added by this notebook.
