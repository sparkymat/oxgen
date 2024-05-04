//nolint:lll,revive
package generator

import (
	"context"
	"fmt"
	"path/filepath"
)

const frontendlistComponentTemplate = `
import {
  Anchor,
  Button,
  Container,
  Flex,
  LoadingOverlay,
  Modal,
  Table,
  Text,
  TextInput,
  Title,
} from '@mantine/core';
import React, {
  ChangeEvent,
  useCallback,
  useEffect,
  useMemo,
  useState,
} from 'react';
import { useParams } from 'react-router-dom';

import { useCreateMutation, useSearchQuery } from '../../slices/{{ .Resource.CamelcaseSingular }}';
import { FilterBar } from '../FilterBar';
import { Pagination } from '../Pagination';
import { {{ .Resource.CamelcaseSingular }} } from '../../models/{{ .Resource.CamelcaseSingular }}';

export const {{ .Resource.CamelcasePlural }}Page = () => {
  const { page: pageString, query } = useParams();

  const [newQuery, setNewQuery] = useState<string>(query || '');
  const [createName, setCreateName] = useState<string>('');
  const [createShown, setCreateShown] = useState<boolean>(false);

  const pageNumber = useMemo((): number => {
    let page = parseInt(pageString || '', 10);

    if (Number.isNaN(page)) {
      page = 1;
    }

    return page;
  }, [pageString]);

  const pageSize = 20;

  const currentFilterURL = useMemo(
    () => (query ?  ` + "`" + `/#/{{ .Resource.UnderscorePlural }}/search/${query}` + "`" + ` : '/#/{{ .Resource.UnderscorePlural }}'),
    [query],
  );

  const newFilterURL = useMemo(
    () => (newQuery ? ` + "`" + `/#/{{ .Resource.UnderscorePlural }}/search/${newQuery}` + "`" + ` : '/#/{{ .Resource.UnderscorePlural }}'),
    [newQuery],
  );

  const { data: itemsData, isLoading: itemsLoading } = useSearchQuery({
    query: query || '',
    pageNumber,
    pageSize,
  });

  const [createItem, { isLoading: isCreating }] = useCreateMutation();

  const createClicked = useCallback(() => {
    createItem({
      name: createName,
    }).then(res => {
      window.location.href = ` + "`" + `/#/{{ .Resource.UnderscorePlural }}/${(res as any).data.id}` + "`" + `;
    });
    setCreateShown(false);
  }, [createItem, createName]);

  const items = useMemo(
    () =>
      itemsData?.items
        ? itemsData?.items.map(i => new {{ .Resource.CamelcaseSingular }}(i))
        : ([] as {{ .Resource.CamelcaseSingular }}[]),
    [itemsData?.items],
  );

  useEffect(() => {
    document.title = '{{ .Resource.CamelcasePlural }}';

    if (query) {
      document.title = ` + "`" + `{{ .Resource.CamelcasePlural }} | ${query}` + "`" + `;
    }
  }, [query]);

  const totalCount = useMemo(
    () => (itemsData ? itemsData.totalCount : 0),
    [itemsData],
  );

  const pageCount = useMemo(
    () => Math.ceil(totalCount / pageSize),
    [totalCount],
  );

  const newQueryChanged = useCallback((evt: ChangeEvent<HTMLInputElement>) => {
    setNewQuery(evt.target.value);
  }, []);

  const createOpened = useCallback(() => {
    setCreateShown(true);
  }, []);

  const createClosed = useCallback(() => {
    setCreateShown(false);
  }, []);

  const createNameChanged = useCallback(
    (evt: ChangeEvent<HTMLInputElement>) => {
      setCreateName(evt.target.value);
    },
    [],
  );

  return (
    <Container fluid>
      <Flex direction="column" gap="md">
        <Title order={3}>{{ .Resource.CamelcasePlural }}</Title>
        <FilterBar
          query={newQuery}
          onQueryChanged={newQueryChanged}
          filterLocation={newFilterURL}
          showCreateModal={createOpened}
        />
        {query && <Text fs="italic">{` + "`" + `Filtering by: ${query}` + "`" + `}</Text>}
        <Table striped highlightOnHover>
          <Table.Tbody>
            {items.map(e => (
              <Table.Tr key={e.id}>
                <Table.Td>
                  <Anchor href={` + "`" + `/#/{{ .Resource.UnderscoreSingular }}/${e.id}` + "`" + `}>
                    <Text size="lg">{e.name}</Text>
                  </Anchor>
                </Table.Td>
              </Table.Tr>
            ))}
          </Table.Tbody>
        </Table>
        <Flex justify="center" mb="md">
          <Pagination
            pageNumber={pageNumber}
            pageCount={pageCount}
            filterURL={currentFilterURL}
          />
        </Flex>
      </Flex>
      <Modal title="New {{ .Resource.CamelcasePlural }}" opened={createShown} onClose={createClosed}>
        <Flex direction="column" gap="md">
          <TextInput
            placeholder="Name"
            value={createName}
            onChange={createNameChanged}
          />
          <Button variant="filled" onClick={createClicked}>
            Add
          </Button>
        </Flex>
        <LoadingOverlay visible={isCreating} />
      </Modal>
      <LoadingOverlay visible={itemsLoading} />
    </Container>
  );
};
`

func (s *Service) generateFrontendComponents(ctx context.Context, input Input) error {
	folderPath := filepath.Join(input.WorkspaceFolder, "frontend", "src", "components", input.Resource.CamelcasePlural()+"Page")

	if err := s.ensureFolderExists(folderPath); err != nil {
		return fmt.Errorf("failed to ensure frontend list component folder exists: %w", err)
	}

	filename := "index.tsx"

	filePath := filepath.Join(folderPath, filename)
	if err := s.appendTemplateToFile(
		ctx,
		filePath,
		0,
		"",
		"frontendList",
		frontendlistComponentTemplate,
		input,
	); err != nil {
		return fmt.Errorf("failed to generate frontend slice: %w", err)
	}

	return nil
}
