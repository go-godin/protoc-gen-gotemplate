// Code generated by protoc-gen-gotemplate
package company

import (
	"github.com/moul/protoc-gen-gotemplate/examples/import/output/models/article"
    "github.com/moul/protoc-gen-gotemplate/examples/import/output/models/common"
)

type Repository interface {
     GetArticle(getarticle *common.GetArticle ) (*company.Article, []*company.Storage,  error)
}







// ------------------------- Public SDK -----------------------------







// GetArticle : proto: missing extension proto: missing extension
func (sdk *Sdk) GetArticle(ctx context.Context, 
  getarticle *GetArticle.GetArticle, token, requestID string)(article *Article.Article, storages []*GetArticleResponse_Storage.Storage, err error) {

  out := &pb.GetArticleResponse{}
	_ = out


  return out.Article, out.Storages, nil
 
}

 
 