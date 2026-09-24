import re
import sys
from collections import Counter
from pathlib import Path

repo = Path(__file__).resolve().parents[4]
sys.path.insert(0, str(repo / 'tools/revkit/scripts'))
from tree3 import parse_tree

tree_paths = sorted((repo / 'data-paul/M16/ttsdata/tree3').rglob('*.tree3'))
trees = [parse_tree(path) for path in tree_paths]
by_shape = {(len(tree.nodes), tree.output_width): tree for tree in trees}
assert len(by_shape) == len(trees), 'tree node-count/output-width shapes are not unique'
lines = (repo / 'tools/revkit/work/stage5/gdb-tree-run.log').read_text().splitlines()

def numbers(line):
    return [int(item) for item in re.findall(r'-?\d+', line.split(':', 1)[1])]

compared = Counter()
for index, line in enumerate(lines):
    scalar = re.search(r'SCALAR_ENTRY .* nodes=(\d+) output_width=(\d+) features=', line)
    vector = re.search(r'VECTOR_ENTRY .* nodes=(\d+) output_width=(\d+) features=', line)
    match = scalar or vector
    if match is None:
        continue
    shape = (int(match.group(1)), int(match.group(2)))
    tree = by_shape.get(shape)
    if tree is None:
        continue
    features = numbers(lines[index + 1]) + numbers(lines[index + 2])
    leaf, expected = tree.evaluate(features)
    if scalar:
        result_line = next(item for item in lines[index + 3 : index + 8]
                           if item.startswith('SCALAR_RETURN'))
        actual = (int(re.search(r'value=(-?\d+)', result_line).group(1)),)
    else:
        result_line = next(item for item in lines[index + 3 : index + 9]
                           if item.startswith('VECTOR_RETURN values:'))
        result_index = lines.index(result_line, index)
        output_line = next(item for item in lines[result_index + 1 :] if item.startswith('0x'))
        actual = tuple(int(item) for item in re.findall(r'-?\d+', result_line.rsplit(':', 1)[1]))
        actual += tuple(numbers(output_line))
    assert actual == expected, (
        f'{tree.path.name}: leaf {leaf} expected {expected}, runtime returned {actual}'
    )
    compared[tree.path.name] += 1

print(f'Runtime tree lookup comparisons: {sum(compared.values())}')
for name, count in sorted(compared.items()):
    print(f'{name}: {count}')
