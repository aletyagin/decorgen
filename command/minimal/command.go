package minimal

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
)

func Run(source string, destination io.Writer, decorName string, packageName string) error {
	fset := token.NewFileSet()
	interfaceNode, err := parser.ParseFile(fset, source, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	if packageName == "" {
		packageName = interfaceNode.Name.Name
	}

	decorNode := &ast.File{
		Name:  ast.NewIdent(packageName),
		Decls: make([]ast.Decl, 0),
	}

	importDecl := &ast.GenDecl{
		Tok:   token.IMPORT,
		Specs: make([]ast.Spec, 0),
	}

	for _, f := range interfaceNode.Decls {
		genD, ok := f.(*ast.GenDecl)
		if !ok {
			fmt.Printf("SKIP %T is not *ast.GenDecl\n", f)
			continue
		}

		for _, spec := range genD.Specs {
			currImport, ok := spec.(*ast.ImportSpec)
			if ok {
				importDecl.Specs = append(importDecl.Specs, currImport)
				continue
			}
		}
	}

	decorNode.Decls = append(decorNode.Decls, importDecl)

	for _, f := range interfaceNode.Decls {
		genD, ok := f.(*ast.GenDecl)
		if !ok {
			fmt.Printf("SKIP %T is not *ast.GenDecl\n", f)
			continue
		}

		for _, spec := range genD.Specs {
			currType, ok := spec.(*ast.TypeSpec)
			if !ok {
				fmt.Printf("SKIP %T is not ast.TypeSpec\n", spec)
				continue
			}

			currInterface, ok := currType.Type.(*ast.InterfaceType)
			if !ok {
				fmt.Printf("SKIP %T is not ast.InterfaceType\n", currInterface)
				continue
			}

			decorNode.Decls = append(
				decorNode.Decls,
				&ast.GenDecl{
					Tok: token.TYPE,
					Specs: []ast.Spec{
						&ast.TypeSpec{
							Name:       ast.NewIdent(decorName),
							TypeParams: currType.TypeParams,
							Type: &ast.StructType{
								Fields: &ast.FieldList{
									List: []*ast.Field{
										{
											Names: []*ast.Ident{{Name: "origin"}},
											Type:  parameterizedType(currType.Name.Name, currType),
										},
									},
								},
							},
						},
					},
				},
			)

			if currInterface.Methods != nil {
				for _, method := range currInterface.Methods.List {

					ast.Inspect(method, func(node ast.Node) bool {
						ftype, ok := method.Type.(*ast.FuncType)
						if !ok {
							return false
						}

						if ftype.Params != nil {
							for f, field := range ftype.Params.List {
								if len(field.Names) == 0 {
									field.Names = []*ast.Ident{ast.NewIdent(fmt.Sprintf("var%d", f))}
								}
							}
						}

						return false
					})

					decorNode.Decls = append(
						decorNode.Decls,
						&ast.FuncDecl{
							Doc: nil,
							Recv: &ast.FieldList{ // receiver (e.g. `(d *RecvDecor[A, B])`)
								List: []*ast.Field{
									{
										Names: []*ast.Ident{{Name: "d"}},
										Type: &ast.StarExpr{
											X: parameterizedType(decorName, currType),
										},
									},
								},
							},
							Name: method.Names[0],
							Type: method.Type.(*ast.FuncType),
							Body: funcBody(method),
						},
					)
				}
			}
		}
	}

	return printer.Fprint(destination, token.NewFileSet(), decorNode)
}

func parameterizedType(typeName string, currType *ast.TypeSpec) ast.Expr {
	if currType.TypeParams != nil {
		types := make([]ast.Expr, 0)
		for _, field := range currType.TypeParams.List {
			for _, name := range field.Names {
				types = append(
					types,
					&ast.Ident{Name: name.Name},
				)
			}
		}

		if len(types) == 1 {
			return &ast.IndexExpr{
				X:     &ast.Ident{Name: typeName},
				Index: types[0],
			}
		}

		return &ast.IndexListExpr{
			X:       &ast.Ident{Name: typeName},
			Indices: types,
		}
	}

	return &ast.Ident{Name: typeName}
}

func funcBody(method *ast.Field) *ast.BlockStmt {
	argNames := make([]ast.Expr, 0)

	ftype, ok := method.Type.(*ast.FuncType)
	if ok {
		if ftype.Params != nil {
			for _, field := range ftype.Params.List {
				for _, name := range field.Names {
					argNames = append(argNames, name)
				}
			}
		}
	}

	call := &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X: &ast.SelectorExpr{
				X: &ast.Ident{Name: "d"},
				Sel: &ast.Ident{
					Name: "origin",
					Obj: &ast.Object{
						Kind: ast.Var,
						Name: "d",
					},
				},
			},
			Sel: &ast.Ident{Name: method.Names[0].Name},
		},
		Args: argNames,
	}

	if ftype.Results == nil {
		return &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{
					X: call,
				},
			},
		}
	}

	return &ast.BlockStmt{
		List: []ast.Stmt{
			&ast.ReturnStmt{
				Results: []ast.Expr{call},
			},
		},
	}
}
